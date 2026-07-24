package openaiClient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchWebCallsResponsesEndpoint(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != ResponsesPath {
			t.Fatalf("Expected path %q, but got %q", ResponsesPath, request.URL.Path)
		}
		if request.Header.Get(AuthorizationHeader) != "Bearer test-key" {
			t.Fatalf("Expected bearer authorization header, but got %q", request.Header.Get(AuthorizationHeader))
		}
		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		if !strings.Contains(string(requestBody), `"type":"web_search"`) {
			t.Fatalf("Expected web_search request body, but got %s", requestBody)
		}

		writer.Header().Set(ContentTypeHeader, "application/json")
		_, _ = writer.Write([]byte(`{"status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"Search worked.","annotations":[]}]}]}`))
	}))
	defer testServer.Close()

	client := Client{APIKey: "test-key", BaseURL: testServer.URL, HTTPClient: testServer.Client()}
	result, err := client.SearchWeb(context.Background(), &WebSearchRequest{
		Model:             "gpt-5.6",
		Query:             "test query",
		SearchContextSize: "low",
	})
	if err != nil {
		t.Fatalf("SearchWeb returned error: %v", err)
	}
	if result.Text != "Search worked." {
		t.Fatalf("Expected search output text, but got %q", result.Text)
	}
}
