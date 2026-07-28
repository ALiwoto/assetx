package appRunner

import (
	"assetx/src/providers/openaiClient"
	"fmt"
	"io"
	"strings"
	"time"
)

func normalizeSearchRequest(request *SearchRequest) error {
	if request == nil {
		return fmt.Errorf("cannot normalize search request because request is nil")
	}

	request.Model = strings.TrimSpace(request.Model)
	if request.Model == "" {
		request.Model = DefaultSearchModel
	}

	request.Query = strings.TrimSpace(request.Query)
	if request.Query == "" {
		return fmt.Errorf("missing required --query value")
	}

	request.SearchContextSize = strings.ToLower(strings.TrimSpace(request.SearchContextSize))
	if request.SearchContextSize == "" {
		request.SearchContextSize = DefaultSearchContextSize
	}
	if !isSupportedSearchContextSize(request.SearchContextSize) {
		return fmt.Errorf(
			"unsupported --context %q: expected %q, %q, or %q",
			request.SearchContextSize,
			SearchContextLow,
			SearchContextMedium,
			SearchContextHigh,
		)
	}

	if request.Timeout == 0 {
		request.Timeout = DefaultSearchTimeout
	}
	if request.Timeout < 0 {
		return fmt.Errorf("invalid --timeout %q: expected a positive duration", request.Timeout)
	}
	if request.ProgressInterval < 0 {
		return fmt.Errorf("invalid --progress-interval %q: expected zero or a positive duration", request.ProgressInterval)
	}

	if len(request.AllowedDomains) > 100 {
		return fmt.Errorf("expected no more than 100 --domain values, but got %d", len(request.AllowedDomains))
	}

	normalizedDomains := make([]string, 0, len(request.AllowedDomains))
	seenDomains := make(map[string]struct{}, len(request.AllowedDomains))
	for _, domain := range request.AllowedDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			return fmt.Errorf("--domain values cannot be empty")
		}
		if strings.Contains(domain, "://") || strings.ContainsAny(domain, "/?#:@ \t\r\n") {
			return fmt.Errorf("invalid --domain %q: use a hostname without a scheme, path, port, query, or fragment", domain)
		}
		if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || !strings.Contains(domain, ".") {
			return fmt.Errorf("invalid --domain %q: expected a hostname such as fab.com", domain)
		}
		if _, alreadySeen := seenDomains[domain]; alreadySeen {
			continue
		}
		seenDomains[domain] = struct{}{}
		normalizedDomains = append(normalizedDomains, domain)
	}
	request.AllowedDomains = normalizedDomains

	return nil
}

func startSearchProgress(writer io.Writer, interval time.Duration, timeout time.Duration) func() {
	if writer == nil || interval <= 0 {
		return func() {}
	}

	startedAt := time.Now()
	_, _ = fmt.Fprintf(
		writer,
		"[assetx]: [progress]: waiting for API response (timeout: %s)...\n",
		timeout,
	)

	done := make(chan struct{})
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				elapsed := time.Since(startedAt).Round(time.Second)
				_, _ = fmt.Fprintf(
					writer,
					"[assetx]: [progress]: still waiting for API response (%s elapsed)...\n",
					elapsed,
				)
			}
		}
	}()

	return func() {
		close(done)
		<-finished
	}
}

func isSupportedSearchContextSize(contextSize string) bool {
	return contextSize == SearchContextLow || contextSize == SearchContextMedium || contextSize == SearchContextHigh
}

func formatWebSearchResult(result openaiClient.WebSearchResult) string {
	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(result.Text))
	if len(result.Sources) == 0 {
		return builder.String()
	}

	builder.WriteString("\n\nSources:\n")
	for _, source := range result.Sources {
		title := strings.TrimSpace(source.Title)
		if title == "" {
			title = source.URL
		}
		title = strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]").Replace(title)
		fmt.Fprintf(&builder, "- [%s](%s)\n", title, source.URL)
	}
	return strings.TrimRight(builder.String(), "\n")
}
