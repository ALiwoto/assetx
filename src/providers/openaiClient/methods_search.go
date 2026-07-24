package openaiClient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (client Client) SearchWeb(ctx context.Context, searchRequest *WebSearchRequest) (WebSearchResult, error) {
	requestBytes, err := buildWebSearchJSONBody(searchRequest)
	if err != nil {
		return WebSearchResult{}, err
	}

	endpointURL := strings.TrimRight(client.BaseURL, "/") + ResponsesPath
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewReader(requestBytes))
	if err != nil {
		return WebSearchResult{}, fmt.Errorf("failed to create Responses API web search request: %w", err)
	}
	httpRequest.Header.Set(AuthorizationHeader, "Bearer "+client.APIKey)
	httpRequest.Header.Set(ContentTypeHeader, "application/json")

	response, err := client.httpClient().Do(httpRequest)
	if err != nil {
		return WebSearchResult{}, fmt.Errorf("Responses API web search request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return WebSearchResult{}, fmt.Errorf("failed to read Responses API web search response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return WebSearchResult{}, parseResponsesAPIError(response.StatusCode, endpointURL, responseBody)
	}

	return decodeWebSearchResponse(responseBody)
}
