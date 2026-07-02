package appRunner

import (
	"assetx/src/core/appConfig"
	"assetx/src/core/imageProcessing"
	"assetx/src/core/openaiClient"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RunImage(request ImageRequest, stdout io.Writer) error {
	normalizedRequest, outputFormat, needsChromaTransparency, err := normalizeImageRequest(request)
	if err != nil {
		return err
	}

	config, err := appConfig.LoadConfig(normalizedRequest.ConfigPath)
	if err != nil {
		return err
	}

	client := openaiClient.NewClient(config)
	imageBytes, err := client.CreateImage(context.Background(), openaiClient.ImageRequest{
		Background:   normalizedRequest.Background,
		Examples:     normalizedRequest.Examples,
		Model:        normalizedRequest.Model,
		OutputFormat: outputFormat,
		Prompt:       normalizedRequest.Prompt,
		Quality:      normalizedRequest.Quality,
		Size:         normalizedRequest.Size,
	})
	if err != nil {
		return err
	}

	if needsChromaTransparency {
		imageBytes, err = imageProcessing.ApplyChromaTransparency(imageBytes)
		if err != nil {
			return err
		}
	}

	if err := ensureOutputDirectory(normalizedRequest.OutputPath); err != nil {
		return err
	}
	if err := os.WriteFile(normalizedRequest.OutputPath, imageBytes, 0644); err != nil {
		return fmt.Errorf("failed to write image output %q: %w", normalizedRequest.OutputPath, err)
	}

	_, _ = fmt.Fprintf(stdout, "Wrote %s\n", normalizedRequest.OutputPath)
	return nil
}

func normalizeImageRequest(request ImageRequest) (ImageRequest, string, bool, error) {
	request.Model = strings.TrimSpace(request.Model)
	if request.Model == "" {
		request.Model = DefaultImageModel
	}
	if !isSupportedImageModel(request.Model) {
		return ImageRequest{}, "", false, fmt.Errorf("unsupported image model %q: supported models are %q and %q", request.Model, ModelGPTImage2, ModelGPTImage15)
	}

	request.Prompt = strings.TrimSpace(request.Prompt)
	if request.Prompt == "" {
		return ImageRequest{}, "", false, fmt.Errorf("missing required --prompt value")
	}

	request.OutputPath = strings.TrimSpace(request.OutputPath)
	if request.OutputPath == "" {
		return ImageRequest{}, "", false, fmt.Errorf("missing required --out value")
	}

	request.Background = strings.ToLower(strings.TrimSpace(request.Background))
	if request.Background == "" {
		request.Background = BackgroundAuto
	}
	if !isSupportedBackground(request.Background) {
		return ImageRequest{}, "", false, fmt.Errorf("unsupported --background %q: expected %q, %q, or %q", request.Background, BackgroundAuto, BackgroundOpaque, BackgroundTransparent)
	}

	request.Quality = strings.ToLower(strings.TrimSpace(request.Quality))
	if request.Quality == "" {
		request.Quality = DefaultImageQuality
	}
	if !isSupportedQuality(request.Quality) {
		return ImageRequest{}, "", false, fmt.Errorf("unsupported --quality %q: expected auto, low, medium, or high", request.Quality)
	}

	request.Size = strings.ToLower(strings.TrimSpace(request.Size))
	if request.Size == "" {
		request.Size = DefaultImageSize
	}
	if err := validateImageSize(request.Model, request.Size); err != nil {
		return ImageRequest{}, "", false, err
	}

	outputFormat, err := outputFormatFromPath(request.OutputPath)
	if err != nil {
		return ImageRequest{}, "", false, err
	}

	needsChromaTransparency := request.Model == ModelGPTImage2 && request.Background == BackgroundTransparent
	if needsChromaTransparency {
		if outputFormat != OutputFormatPNG {
			return ImageRequest{}, "", false, fmt.Errorf("transparent output with %s requires a .png --out path because the chroma transparency post-processor writes PNG alpha", ModelGPTImage2)
		}
		request.Background = BackgroundOpaque
		request.Prompt = appendChromaPrompt(request.Prompt)
	}

	for _, examplePath := range request.Examples {
		if strings.TrimSpace(examplePath) == "" {
			return ImageRequest{}, "", false, fmt.Errorf("--example values cannot be empty")
		}
		if _, err := os.Stat(examplePath); err != nil {
			return ImageRequest{}, "", false, fmt.Errorf("failed to access --example %q: %w", examplePath, err)
		}
	}

	return request, outputFormat, needsChromaTransparency, nil
}

func isSupportedImageModel(model string) bool {
	return model == ModelGPTImage2 || model == ModelGPTImage15
}

func isSupportedBackground(background string) bool {
	return background == BackgroundAuto || background == BackgroundOpaque || background == BackgroundTransparent
}

func isSupportedQuality(quality string) bool {
	return quality == "auto" || quality == "low" || quality == "medium" || quality == "high"
}

func validateImageSize(model string, size string) error {
	if size == "auto" {
		return nil
	}

	parts := strings.Split(size, "x")
	if len(parts) != 2 {
		return fmt.Errorf("invalid --size %q: expected auto or WIDTHxHEIGHT such as 1024x1024", size)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid --size %q: width must be a number", size)
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid --size %q: height must be a number", size)
	}
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid --size %q: width and height must both be positive", size)
	}

	if model == ModelGPTImage2 {
		if width%16 != 0 || height%16 != 0 {
			return fmt.Errorf("invalid --size %q for %s: width and height must both be divisible by 16", size, ModelGPTImage2)
		}
		if width > 3840 || height > 2160 {
			return fmt.Errorf("invalid --size %q for %s: maximum supported edge is 3840x2160", size, ModelGPTImage2)
		}
		if width > height*3 || height > width*3 {
			return fmt.Errorf("invalid --size %q for %s: aspect ratio must be between 1:3 and 3:1", size, ModelGPTImage2)
		}
		return nil
	}

	if size != "1024x1024" && size != "1024x1536" && size != "1536x1024" {
		return fmt.Errorf("invalid --size %q for %s: expected auto, 1024x1024, 1024x1536, or 1536x1024", size, model)
	}

	return nil
}

func outputFormatFromPath(outputPath string) (string, error) {
	extension := strings.ToLower(filepath.Ext(outputPath))
	switch extension {
	case ".png":
		return OutputFormatPNG, nil
	case ".jpg", ".jpeg":
		return OutputFormatJPEG, nil
	case ".webp":
		return OutputFormatWEBP, nil
	default:
		return "", fmt.Errorf("unsupported --out extension %q: expected .png, .jpg, .jpeg, or .webp", extension)
	}
}

func appendChromaPrompt(prompt string) string {
	return prompt + "\n\nRender the asset on a perfectly flat #00FF00 chroma key background. Do not use shadows, glow, reflections, gradients, texture, or green details touching the background. Keep the asset fully inside the canvas."
}

func ensureOutputDirectory(outputPath string) error {
	outputDirectory := filepath.Dir(outputPath)
	if outputDirectory == "." || outputDirectory == "" {
		return nil
	}
	if err := os.MkdirAll(outputDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outputDirectory, err)
	}
	return nil
}
