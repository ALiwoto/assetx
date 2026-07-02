package openaiClient

import (
	"assetx/src/core/appConfig"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func NewClient(config *appConfig.AssetxConfig) Client {
	return Client{
		APIKey:  strings.TrimSpace(config.APIKey),
		BaseURL: normalizeBaseURL(config.ProxyBaseURL),
	}
}

func imageEndpoint(baseURL string, endpointPath string) string {
	return strings.TrimRight(baseURL, "/") + endpointPath
}

func normalizeBaseURL(baseURL string) string {
	trimmedBaseURL := strings.TrimSpace(baseURL)
	if trimmedBaseURL == "" {
		return DefaultOpenAIBaseURL
	}

	parsedURL, err := url.Parse(trimmedBaseURL)
	if err != nil {
		return strings.TrimRight(trimmedBaseURL, "/")
	}

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		if strings.EqualFold(parsedURL.Host, "api.openai.com") {
			parsedURL.Path = "/v1"
			return strings.TrimRight(parsedURL.String(), "/")
		}
		parsedURL.Path = "/openai/v1"
		return strings.TrimRight(parsedURL.String(), "/")
	}

	return strings.TrimRight(trimmedBaseURL, "/")
}

func buildImageRequestBody(imageRequest ImageRequest) (io.Reader, string, string, func(), error) {
	if len(imageRequest.Examples) == 0 {
		requestBytes, err := buildGenerationJSONBody(imageRequest)
		if err != nil {
			return nil, "", "", func() {}, err
		}
		return bytes.NewReader(requestBytes), "application/json", ImageGenerationsPath, func() {}, nil
	}

	requestBody, contentType, cleanup, err := buildEditMultipartBody(imageRequest)
	if err != nil {
		return nil, "", "", func() {}, err
	}
	return requestBody, contentType, ImageEditsPath, cleanup, nil
}

func buildGenerationJSONBody(imageRequest ImageRequest) ([]byte, error) {
	requestMap := map[string]string{
		"model":         imageRequest.Model,
		"output_format": imageRequest.OutputFormat,
		"prompt":        imageRequest.Prompt,
		"quality":       imageRequest.Quality,
		"size":          imageRequest.Size,
	}
	if imageRequest.Background != "" && imageRequest.Background != "auto" {
		requestMap["background"] = imageRequest.Background
	}

	requestBytes, err := json.Marshal(requestMap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image generation request: %w", err)
	}
	return requestBytes, nil
}

func buildEditMultipartBody(imageRequest ImageRequest) (io.Reader, string, func(), error) {
	bodyBuffer := bytes.NewBuffer(nil)
	multipartWriter := multipart.NewWriter(bodyBuffer)
	openFiles := make([]*os.File, 0, len(imageRequest.Examples))

	cleanup := func() {
		for _, file := range openFiles {
			_ = file.Close()
		}
	}

	fields := map[string]string{
		"model":         imageRequest.Model,
		"output_format": imageRequest.OutputFormat,
		"prompt":        imageRequest.Prompt,
		"quality":       imageRequest.Quality,
		"size":          imageRequest.Size,
	}
	if imageRequest.Background != "" && imageRequest.Background != "auto" {
		fields["background"] = imageRequest.Background
	}
	for fieldName, fieldValue := range fields {
		if err := multipartWriter.WriteField(fieldName, fieldValue); err != nil {
			cleanup()
			return nil, "", func() {}, fmt.Errorf("failed to write multipart field %q: %w", fieldName, err)
		}
	}

	for index, examplePath := range imageRequest.Examples {
		exampleFile, err := os.Open(examplePath)
		if err != nil {
			cleanup()
			return nil, "", func() {}, fmt.Errorf("failed to open example image %q: %w", examplePath, err)
		}
		openFiles = append(openFiles, exampleFile)

		part, err := multipartWriter.CreatePart(exampleFilePartHeader(examplePath, index))
		if err != nil {
			cleanup()
			return nil, "", func() {}, fmt.Errorf("failed to create multipart image field for %q: %w", examplePath, err)
		}
		if _, err := io.Copy(part, exampleFile); err != nil {
			cleanup()
			return nil, "", func() {}, fmt.Errorf("failed to read example image %q into request: %w", examplePath, err)
		}
	}

	if err := multipartWriter.Close(); err != nil {
		cleanup()
		return nil, "", func() {}, fmt.Errorf("failed to finalize multipart image request: %w", err)
	}

	return bodyBuffer, multipartWriter.FormDataContentType(), cleanup, nil
}

func exampleFilePartHeader(examplePath string, index int) textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image[]"; filename="%s"`, multipartEscape(filepath.Base(examplePath))))
	header.Set("Content-Type", mimeTypeFromExtension(examplePath))
	header.Set("X-Assetx-Example-Index", fmt.Sprintf("%d", index))
	return header
}

func mimeTypeFromExtension(filePath string) string {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func multipartEscape(value string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(value)
}

func decodeImageResponse(responseBody []byte) (ImageResponse, error) {
	var imageResponse ImageResponse
	if err := json.Unmarshal(responseBody, &imageResponse); err != nil {
		return ImageResponse{}, fmt.Errorf("failed to parse image API response as JSON: %w", err)
	}
	return imageResponse, nil
}

func decodeBase64Image(encodedImage string) ([]byte, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(encodedImage)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image data from API response: %w", err)
	}
	return imageBytes, nil
}

func parseAPIError(statusCode int, endpointURL string, responseBody []byte) error {
	var errorResponse APIErrorResponse
	if err := json.Unmarshal(responseBody, &errorResponse); err == nil && errorResponse.Error.Message != "" {
		if statusCode == 404 {
			return fmt.Errorf("image API endpoint %q returned HTTP 404: %s. Verify that proxy_base_url points to an OpenAI-compatible /v1 API that supports /images/generations and /images/edits", endpointURL, errorResponse.Error.Message)
		}
		return fmt.Errorf("image API endpoint %q returned HTTP %d: %s", endpointURL, statusCode, errorResponse.Error.Message)
	}
	return fmt.Errorf("image API endpoint %q returned HTTP %d: %s", endpointURL, statusCode, string(responseBody))
}
