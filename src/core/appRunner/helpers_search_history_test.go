package appRunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteSearchHistoryCreatesTimestampedMarkdown(t *testing.T) {
	homeDir := t.TempDir()
	searchedAt := time.Date(2026, time.July, 28, 17, 5, 3, 123456789, time.UTC)
	request := &SearchRequest{
		AllowedDomains:    []string{"github.com", "pkg.go.dev"},
		Model:             "gpt-5.6",
		Query:             "find bbolt documentation",
		SearchContextSize: "medium",
	}

	historyPath, err := writeSearchHistory(homeDir, request, "Search result.\n\nSources:\n- [Docs](https://example.com)", searchedAt)
	if err != nil {
		t.Fatalf("writeSearchHistory returned error: %v", err)
	}

	expectedPath := filepath.Join(
		homeDir,
		".assetx",
		"search_history",
		"search_20260728T170503.123456789Z.md",
	)
	if historyPath != expectedPath {
		t.Fatalf("Expected history path %q, but got %q", expectedPath, historyPath)
	}

	historyBytes, err := os.ReadFile(historyPath)
	if err != nil {
		t.Fatalf("Failed to read search history: %v", err)
	}
	history := string(historyBytes)
	for _, expected := range []string{
		"# assetx search history",
		"find bbolt documentation",
		"Search result.",
		"`github.com`, `pkg.go.dev`",
	} {
		if !strings.Contains(history, expected) {
			t.Fatalf("Expected history to contain %q, but got %q", expected, history)
		}
	}
}

func TestWriteSearchHistoryDoesNotOverwriteTimestampCollision(t *testing.T) {
	homeDir := t.TempDir()
	searchedAt := time.Date(2026, time.July, 28, 17, 5, 3, 0, time.UTC)
	request := &SearchRequest{Query: "test"}

	if _, err := writeSearchHistory(homeDir, request, "first", searchedAt); err != nil {
		t.Fatalf("First writeSearchHistory returned error: %v", err)
	}
	if _, err := writeSearchHistory(homeDir, request, "second", searchedAt); err == nil {
		t.Fatal("Expected timestamp collision to fail instead of overwriting history")
	}
}
