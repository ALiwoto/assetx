package appRunner

import (
	"assetx/src/core/appConfig"
	"assetx/src/providers/openaiClient"
	"context"
	"errors"
	"fmt"
	"io"
)

func RunSearch(request *SearchRequest, stdout io.Writer, stderr io.Writer) error {
	if err := normalizeSearchRequest(request); err != nil {
		return err
	}

	config, err := appConfig.LoadConfig(request.ConfigPath)
	if err != nil {
		return err
	}

	client := openaiClient.NewClient(config)
	ctx, cancel := context.WithTimeout(context.Background(), request.Timeout)
	defer cancel()

	stopProgress := startSearchProgress(stderr, request.ProgressInterval, request.Timeout)
	result, err := client.SearchWeb(ctx, &openaiClient.WebSearchRequest{
		AllowedDomains:    request.AllowedDomains,
		Model:             request.Model,
		Query:             request.Query,
		SearchContextSize: request.SearchContextSize,
	})
	stopProgress()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf(
				"web search timed out after %s while waiting for the Responses API; increase --timeout and ensure the calling process has a longer runtime limit",
				request.Timeout,
			)
		}
		return err
	}

	_, _ = fmt.Fprintln(stdout, formatWebSearchResult(result))
	return nil
}
