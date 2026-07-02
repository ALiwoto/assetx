package openaiClient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client Client) CreateImage(ctx context.Context, imageRequest ImageRequest) ([]byte, error) {
	requestBody, contentType, endpointPath, cleanup, err := buildImageRequestBody(imageRequest)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, imageEndpoint(client.BaseURL, endpointPath), requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create image API request: %w", err)
	}
	httpRequest.Header.Set(AuthorizationHeader, "Bearer "+client.APIKey)
	httpRequest.Header.Set(ContentTypeHeader, contentType)

	response, err := client.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("image API request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image API response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, parseAPIError(response.StatusCode, imageEndpoint(client.BaseURL, endpointPath), responseBody)
	}

	imageResponse, err := decodeImageResponse(responseBody)
	if err != nil {
		return nil, err
	}

	return client.extractImageBytes(ctx, imageResponse)
}

func (client Client) extractImageBytes(ctx context.Context, imageResponse ImageResponse) ([]byte, error) {
	if len(imageResponse.Data) != 1 {
		return nil, fmt.Errorf("expected image API response data length of 1, but got %d", len(imageResponse.Data))
	}

	imageData := imageResponse.Data[0]
	if imageData.B64JSON != "" {
		return decodeBase64Image(imageData.B64JSON)
	}
	if imageData.URL != "" {
		return client.downloadImageURL(ctx, imageData.URL)
	}

	return nil, fmt.Errorf("image API response did not include b64_json or url")
}

func (client Client) downloadImageURL(ctx context.Context, imageURL string) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, bytes.NewReader(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to create image download request: %w", err)
	}

	response, err := client.httpClient().Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to download generated image from URL: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read generated image download response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("generated image download failed with HTTP %d: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func (client Client) httpClient() *http.Client {
	if client.HTTPClient != nil {
		return client.HTTPClient
	}
	return http.DefaultClient
}
