package appRunner

import (
	"assetx/src/providers/openaiClient"
	"strings"
	"testing"
)

func TestNormalizeSearchRequestDefaultsAndDeduplicatesDomains(t *testing.T) {
	request := &SearchRequest{
		AllowedDomains: []string{" FAB.com ", "fab.com"},
		Query:          " find characters ",
	}

	if err := normalizeSearchRequest(request); err != nil {
		t.Fatalf("normalizeSearchRequest returned error: %v", err)
	}
	if request.Model != DefaultSearchModel {
		t.Fatalf("Expected default model %q, but got %q", DefaultSearchModel, request.Model)
	}
	if request.SearchContextSize != DefaultSearchContextSize {
		t.Fatalf("Expected default context %q, but got %q", DefaultSearchContextSize, request.SearchContextSize)
	}
	if len(request.AllowedDomains) != 1 || request.AllowedDomains[0] != "fab.com" {
		t.Fatalf("Expected one normalized fab.com domain, but got %+v", request.AllowedDomains)
	}
}

func TestNormalizeSearchRequestRejectsDomainURL(t *testing.T) {
	err := normalizeSearchRequest(&SearchRequest{
		AllowedDomains: []string{"https://fab.com/search"},
		Query:          "find characters",
	})
	if err == nil || !strings.Contains(err.Error(), "without a scheme") {
		t.Fatalf("Expected domain URL validation error, but got %v", err)
	}
}

func TestFormatWebSearchResultIncludesClickableSources(t *testing.T) {
	formatted := formatWebSearchResult(openaiClient.WebSearchResult{
		Text: "Found a result.",
		Sources: []openaiClient.WebSearchSource{
			{Title: "A [modular] pack", URL: "https://www.fab.com/listings/example"},
		},
	})
	if !strings.Contains(formatted, `[A \[modular\] pack](https://www.fab.com/listings/example)`) {
		t.Fatalf("Expected clickable escaped source link, but got %q", formatted)
	}
}
