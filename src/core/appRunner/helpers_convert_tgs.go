package appRunner

import (
	"assetx/src/core/imageProcessing"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func normalizeConvertTGSRequest(request *ConvertTGSRequest) error {
	if request == nil {
		return fmt.Errorf("cannot normalize convert-tgs request because request is nil")
	}

	request.FFMPEGPath = strings.TrimSpace(request.FFMPEGPath)
	if request.FFMPEGPath == "" {
		request.FFMPEGPath = imageProcessing.DefaultFFMPEGExecutable
	}

	request.InputPath = strings.TrimSpace(request.InputPath)
	if request.InputPath == "" {
		return fmt.Errorf("missing required --in value")
	}
	if _, err := os.Stat(request.InputPath); err != nil {
		return fmt.Errorf("failed to access --in %q: %w", request.InputPath, err)
	}

	request.OutputPath = strings.TrimSpace(request.OutputPath)
	if request.OutputPath == "" {
		return fmt.Errorf("missing required --out value")
	}
	if strings.ToLower(filepath.Ext(request.OutputPath)) != ".png" {
		return fmt.Errorf("convert-tgs requires a .png --out path because the output is a PNG sprite sheet")
	}

	return nil
}
