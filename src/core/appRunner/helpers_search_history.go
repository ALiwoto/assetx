package appRunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func saveSearchHistory(request *SearchRequest, formattedResult string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to locate the user home directory for search history: %w", err)
	}

	_, err = writeSearchHistory(homeDir, request, formattedResult, time.Now())
	return err
}

func writeSearchHistory(
	homeDir string,
	request *SearchRequest,
	formattedResult string,
	searchedAt time.Time,
) (string, error) {
	if strings.TrimSpace(homeDir) == "" {
		return "", fmt.Errorf("cannot save search history because the user home directory is empty")
	}
	if request == nil {
		return "", fmt.Errorf("cannot save search history because the search request is nil")
	}

	historyDir := filepath.Join(homeDir, ".assetx", "search_history")
	if err := os.MkdirAll(historyDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create search history directory %q: %w", historyDir, err)
	}

	timestamp := searchedAt.UTC().Format("20060102T150405.000000000Z")
	historyPath := filepath.Join(historyDir, "search_"+timestamp+".md")
	historyContent := formatSearchHistory(request, formattedResult, searchedAt)

	historyFile, err := os.OpenFile(historyPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", fmt.Errorf("failed to create search history file %q: %w", historyPath, err)
	}

	if _, err := historyFile.WriteString(historyContent); err != nil {
		_ = historyFile.Close()
		_ = os.Remove(historyPath)
		return "", fmt.Errorf("failed to write search history file %q: %w", historyPath, err)
	}
	if err := historyFile.Close(); err != nil {
		_ = os.Remove(historyPath)
		return "", fmt.Errorf("failed to close search history file %q: %w", historyPath, err)
	}

	return historyPath, nil
}

func formatSearchHistory(request *SearchRequest, formattedResult string, searchedAt time.Time) string {
	var builder strings.Builder
	builder.WriteString("# assetx search history\n\n")
	fmt.Fprintf(&builder, "- Searched at: `%s`\n", searchedAt.UTC().Format(time.RFC3339Nano))
	fmt.Fprintf(&builder, "- Model: `%s`\n", request.Model)
	fmt.Fprintf(&builder, "- Search context: `%s`\n", request.SearchContextSize)
	if len(request.AllowedDomains) > 0 {
		fmt.Fprintf(&builder, "- Allowed domains: `%s`\n", strings.Join(request.AllowedDomains, "`, `"))
	}
	builder.WriteString("\n## Query\n\n")
	builder.WriteString(request.Query)
	builder.WriteString("\n\n## Result\n\n")
	builder.WriteString(strings.TrimSpace(formattedResult))
	builder.WriteString("\n")

	return builder.String()
}
