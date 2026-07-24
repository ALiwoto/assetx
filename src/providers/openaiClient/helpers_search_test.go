package openaiClient

import (
	"encoding/json"
	"testing"
)

func TestBuildWebSearchJSONBodyIncludesRequiredToolAndDomains(t *testing.T) {
	requestBody, err := buildWebSearchJSONBody(&WebSearchRequest{
		AllowedDomains:    []string{"fab.com"},
		Model:             "gpt-5.6",
		Query:             "find modular medieval characters",
		SearchContextSize: "high",
	})
	if err != nil {
		t.Fatalf("buildWebSearchJSONBody returned error: %v", err)
	}

	var decodedRequest responsesCreateRequest
	if err := json.Unmarshal(requestBody, &decodedRequest); err != nil {
		t.Fatalf("Failed to decode request JSON: %v", err)
	}
	if decodedRequest.ToolChoice != "required" {
		t.Fatalf("Expected required tool choice, but got %q", decodedRequest.ToolChoice)
	}
	if len(decodedRequest.Tools) != 1 || decodedRequest.Tools[0].Type != "web_search" {
		t.Fatalf("Expected exactly one web_search tool, but got %+v", decodedRequest.Tools)
	}
	if decodedRequest.Tools[0].Filters == nil || len(decodedRequest.Tools[0].Filters.AllowedDomains) != 1 {
		t.Fatalf("Expected one allowed domain, but got %+v", decodedRequest.Tools[0].Filters)
	}
	if decodedRequest.Tools[0].Filters.AllowedDomains[0] != "fab.com" {
		t.Fatalf("Expected fab.com domain, but got %q", decodedRequest.Tools[0].Filters.AllowedDomains[0])
	}
}

func TestDecodeWebSearchResponseCollectsUniqueSources(t *testing.T) {
	responseBody := []byte(`{
		"status":"completed",
		"output":[
			{"type":"web_search_call","action":{"sources":[
				{"type":"url","title":"Fab result","url":"https://www.fab.com/listings/example"}
			]}},
			{"type":"message","content":[{"type":"output_text","text":"Found one result.","annotations":[
				{"type":"url_citation","title":"Fab result","url":"https://www.fab.com/listings/example"},
				{"type":"url_citation","title":"Documentation","url":"https://example.com/docs"}
			]}]}
		]
	}`)

	result, err := decodeWebSearchResponse(responseBody)
	if err != nil {
		t.Fatalf("decodeWebSearchResponse returned error: %v", err)
	}
	if result.Text != "Found one result." {
		t.Fatalf("Expected output text, but got %q", result.Text)
	}
	if len(result.Sources) != 2 {
		t.Fatalf("Expected 2 unique sources, but got %d: %+v", len(result.Sources), result.Sources)
	}
}
