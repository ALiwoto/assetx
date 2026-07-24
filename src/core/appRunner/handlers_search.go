package appRunner

import (
	"assetx/src/core/appConfig"
	"assetx/src/providers/openaiClient"
	"context"
	"fmt"
	"io"
)

func RunSearch(request *SearchRequest, stdout io.Writer) error {
	if err := normalizeSearchRequest(request); err != nil {
		return err
	}

	config, err := appConfig.LoadConfig(request.ConfigPath)
	if err != nil {
		return err
	}

	client := openaiClient.NewClient(config)
	result, err := client.SearchWeb(context.Background(), &openaiClient.WebSearchRequest{
		AllowedDomains:    request.AllowedDomains,
		Model:             request.Model,
		Query:             request.Query,
		SearchContextSize: request.SearchContextSize,
	})
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(stdout, formatWebSearchResult(result))
	return nil
}
