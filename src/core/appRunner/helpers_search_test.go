package appRunner

import (
	"assetx/src/providers/openaiClient"
	"bytes"
	"strings"
	"testing"
	"time"
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
	if request.Timeout != DefaultSearchTimeout {
		t.Fatalf("Expected default timeout %q, but got %q", DefaultSearchTimeout, request.Timeout)
	}
	if len(request.AllowedDomains) != 1 || request.AllowedDomains[0] != "fab.com" {
		t.Fatalf("Expected one normalized fab.com domain, but got %+v", request.AllowedDomains)
	}
}

func TestNormalizeSearchRequestRejectsNegativeDurations(t *testing.T) {
	for fieldName, request := range map[string]*SearchRequest{
		"progress interval": {
			ProgressInterval: -time.Second,
			Query:            "find characters",
		},
		"timeout": {
			Query:   "find characters",
			Timeout: -time.Second,
		},
	} {
		t.Run(fieldName, func(t *testing.T) {
			if err := normalizeSearchRequest(request); err == nil {
				t.Fatalf("Expected negative %s to fail", fieldName)
			}
		})
	}
}

func TestStartSearchProgressReportsWhileWaiting(t *testing.T) {
	var output bytes.Buffer
	stop := startSearchProgress(&output, time.Millisecond, time.Minute)
	time.Sleep(5 * time.Millisecond)
	stop()

	progressOutput := output.String()
	if !strings.Contains(progressOutput, "waiting for API response") {
		t.Fatalf("Expected initial progress output, but got %q", progressOutput)
	}
	if !strings.Contains(progressOutput, "still waiting for API response") {
		t.Fatalf("Expected periodic progress output, but got %q", progressOutput)
	}
}

func TestStartSearchProgressAllowsDisabling(t *testing.T) {
	var output bytes.Buffer
	stop := startSearchProgress(&output, 0, time.Minute)
	stop()

	if output.Len() != 0 {
		t.Fatalf("Expected disabled progress to produce no output, but got %q", output.String())
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
