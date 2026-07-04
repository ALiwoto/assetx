package appRunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func normalizeConvertWEBPRequest(request *ConvertWEBPRequest) error {
	if request == nil {
		return fmt.Errorf("cannot normalize convert-webp request because request is nil")
	}

	request.InputPath = strings.TrimSpace(request.InputPath)
	if request.InputPath == "" {
		return fmt.Errorf("missing required --in value or positional input path")
	}
	if strings.ToLower(filepath.Ext(request.InputPath)) != ".webp" {
		return fmt.Errorf("convert-webp requires a .webp input path")
	}
	if _, err := os.Stat(request.InputPath); err != nil {
		return fmt.Errorf("failed to access --in %q: %w", request.InputPath, err)
	}

	request.OutputPath = strings.TrimSpace(request.OutputPath)
	if request.OutputPath == "" {
		request.OutputPath = request.InputPath + ".png"
	}
	if strings.ToLower(filepath.Ext(request.OutputPath)) != ".png" {
		return fmt.Errorf("convert-webp requires a .png --out path")
	}

	return nil
}
