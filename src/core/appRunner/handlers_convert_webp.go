package appRunner

import (
	"assetx/src/converters/webpConverter"
	"fmt"
	"io"
)

func RunConvertWEBP(request *ConvertWEBPRequest, stdout io.Writer) error {
	if err := normalizeConvertWEBPRequest(request); err != nil {
		return err
	}

	if err := ensureOutputDirectory(request.OutputPath); err != nil {
		return err
	}

	if err := webpConverter.ConvertWebpToPNG(&webpConverter.ConvertWEBPOptions{
		InputPath:  request.InputPath,
		OutputPath: request.OutputPath,
	}); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(stdout, "Wrote %s\n", request.OutputPath)
	return nil
}
