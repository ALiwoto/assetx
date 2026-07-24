package openaiClient

import (
	"encoding/json"
	"fmt"
	"strings"
)

func buildWebSearchJSONBody(searchRequest *WebSearchRequest) ([]byte, error) {
	if searchRequest == nil {
		return nil, fmt.Errorf("cannot build web search request body because request is nil")
	}

	tool := webSearchTool{
		SearchContextSize: searchRequest.SearchContextSize,
		Type:              "web_search",
	}
	if len(searchRequest.AllowedDomains) > 0 {
		tool.Filters = &webSearchFilters{AllowedDomains: searchRequest.AllowedDomains}
	}

	requestBytes, err := json.Marshal(responsesCreateRequest{
		Include:    []string{"web_search_call.action.sources"},
		Input:      searchRequest.Query,
		Model:      searchRequest.Model,
		ToolChoice: "required",
		Tools:      []webSearchTool{tool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode web search request: %w", err)
	}
	return requestBytes, nil
}

func decodeWebSearchResponse(responseBody []byte) (WebSearchResult, error) {
	var apiResponse responsesAPIResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return WebSearchResult{}, fmt.Errorf("failed to parse Responses API web search response as JSON: %w", err)
	}
	if apiResponse.Error != nil {
		return WebSearchResult{}, fmt.Errorf("Responses API web search failed: %s", apiResponse.Error.Message)
	}

	result := WebSearchResult{}
	for _, outputItem := range apiResponse.Output {
		if outputItem.Type == "web_search_call" && outputItem.Action != nil {
			for _, source := range outputItem.Action.Sources {
				appendUniqueWebSearchSource(&result.Sources, source)
			}
		}
		if outputItem.Type != "message" {
			continue
		}
		for _, content := range outputItem.Content {
			if content.Type != "output_text" || strings.TrimSpace(content.Text) == "" {
				continue
			}
			if result.Text != "" {
				result.Text += "\n"
			}
			result.Text += content.Text
			for _, annotation := range content.Annotations {
				if annotation.Type != "url_citation" {
					continue
				}
				appendUniqueWebSearchSource(&result.Sources, WebSearchSource{
					Title: annotation.Title,
					Type:  annotation.Type,
					URL:   annotation.URL,
				})
			}
		}
	}

	if strings.TrimSpace(result.Text) == "" {
		return WebSearchResult{}, fmt.Errorf(
			"Responses API web search returned status %q but no output text",
			apiResponse.Status,
		)
	}
	return result, nil
}

func appendUniqueWebSearchSource(sources *[]WebSearchSource, source WebSearchSource) {
	source.URL = strings.TrimSpace(source.URL)
	if source.URL == "" {
		return
	}
	for _, existingSource := range *sources {
		if existingSource.URL == source.URL {
			return
		}
	}
	*sources = append(*sources, source)
}

func parseResponsesAPIError(statusCode int, endpointURL string, responseBody []byte) error {
	var errorResponse APIErrorResponse
	if err := json.Unmarshal(responseBody, &errorResponse); err == nil && errorResponse.Error.Message != "" {
		if statusCode == 404 {
			return fmt.Errorf(
				"Responses API endpoint %q returned HTTP 404: %s. Verify that proxy_base_url points to an OpenAI-compatible /v1 API that supports /responses",
				endpointURL,
				errorResponse.Error.Message,
			)
		}
		return fmt.Errorf("Responses API endpoint %q returned HTTP %d: %s", endpointURL, statusCode, errorResponse.Error.Message)
	}
	return fmt.Errorf("Responses API endpoint %q returned HTTP %d: %s", endpointURL, statusCode, string(responseBody))
}
