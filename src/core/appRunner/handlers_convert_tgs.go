package appRunner

import (
	"assetx/src/core/imageProcessing"
	"context"
	"fmt"
	"io"
)

func RunConvertTGS(request *ConvertTGSRequest, stdout io.Writer) error {
	if err := normalizeConvertTGSRequest(request); err != nil {
		return err
	}

	if err := ensureOutputDirectory(request.OutputPath); err != nil {
		return err
	}

	if err := imageProcessing.ConvertTGSToSpritePNG(context.Background(), imageProcessing.ConvertTGSOptions{
		FFMPEGPath: request.FFMPEGPath,
		InputPath:  request.InputPath,
		OutputPath: request.OutputPath,
	}); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(stdout, "Wrote %s\n", request.OutputPath)
	return nil
}
