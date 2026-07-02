package openaiClient

import "testing"

func TestNormalizeBaseURLUsesOpenAIDefaultForEmptyValue(t *testing.T) {
	baseURL := normalizeBaseURL("")
	if baseURL != DefaultOpenAIBaseURL {
		t.Fatalf("Expected %q, but got %q", DefaultOpenAIBaseURL, baseURL)
	}
}

func TestNormalizeBaseURLKeepsProxyRootBehavior(t *testing.T) {
	baseURL := normalizeBaseURL("https://my.proxydomain.com/")
	if baseURL != "https://my.proxydomain.com/openai/v1" {
		t.Fatalf("Expected proxy OpenAI-compatible path, but got %q", baseURL)
	}
}

func TestNormalizeBaseURLUsesV1ForOpenAIHostRoot(t *testing.T) {
	baseURL := normalizeBaseURL("https://api.openai.com/")
	if baseURL != DefaultOpenAIBaseURL {
		t.Fatalf("Expected %q, but got %q", DefaultOpenAIBaseURL, baseURL)
	}
}
