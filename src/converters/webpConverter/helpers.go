package webpConverter

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

func ConvertWebpToPNG(options *ConvertWEBPOptions) error {
	if err := normalizeConvertWEBPOptions(options); err != nil {
		return err
	}

	inputFile, err := os.Open(options.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open input WebP %q: %w", options.InputPath, err)
	}
	defer inputFile.Close()

	sourceImage, err := webp.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("failed to decode input WebP %q: %w", options.InputPath, err)
	}
	sourceBounds := sourceImage.Bounds()
	if sourceBounds.Dx() <= 0 || sourceBounds.Dy() <= 0 {
		return fmt.Errorf("expected positive WebP dimensions, but got %dx%d", sourceBounds.Dx(), sourceBounds.Dy())
	}

	outputFile, err := os.Create(options.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create PNG output %q: %w", options.OutputPath, err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, sourceImage); err != nil {
		return fmt.Errorf("failed to encode PNG output %q: %w", options.OutputPath, err)
	}

	return nil
}

func normalizeConvertWEBPOptions(options *ConvertWEBPOptions) error {
	if options == nil {
		return fmt.Errorf("cannot normalize WebP conversion options because options is nil")
	}

	options.InputPath = strings.TrimSpace(options.InputPath)
	if options.InputPath == "" {
		return fmt.Errorf("missing input WebP path")
	}
	if strings.ToLower(filepath.Ext(options.InputPath)) != ".webp" {
		return fmt.Errorf("WebP conversion requires a .webp input path")
	}
	if _, err := os.Stat(options.InputPath); err != nil {
		return fmt.Errorf("failed to access input WebP %q: %w", options.InputPath, err)
	}

	options.OutputPath = strings.TrimSpace(options.OutputPath)
	if options.OutputPath == "" {
		return fmt.Errorf("missing output PNG path")
	}
	if strings.ToLower(filepath.Ext(options.OutputPath)) != ".png" {
		return fmt.Errorf("WebP conversion requires a .png output path")
	}

	return nil
}
