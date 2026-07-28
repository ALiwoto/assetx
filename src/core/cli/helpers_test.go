package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunConvertWEBPHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"help", "convert-webp"}, &stdout, &stderr)
	if exitCode != ExitSuccess {
		t.Fatalf("Expected success exit code, but got %d with stderr %q", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "assetx convert-webp") {
		t.Fatalf("Expected convert-webp help, but got %q", stdout.String())
	}
}

func TestRunConvertWEBPRejectsInputFlagAndPositionalInput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"convert-webp", "--in", "input.webp", "other.webp"}, &stdout, &stderr)
	if exitCode != ExitFailure {
		t.Fatalf("Expected failure exit code, but got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "use either --in or one positional input path") {
		t.Fatalf("Expected duplicate input error, but got %q", stderr.String())
	}
}

func TestRunSearchHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"help", "search"}, &stdout, &stderr)
	if exitCode != ExitSuccess {
		t.Fatalf("Expected success exit code, but got %d with stderr %q", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "assetx search") {
		t.Fatalf("Expected search help, but got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "--timeout duration") {
		t.Fatalf("Expected search timeout help, but got %q", stdout.String())
	}
}

func TestRunSearchRejectsInvalidTimeout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"search", "--query", "test", "--timeout", "later"}, &stdout, &stderr)
	if exitCode != ExitFailure {
		t.Fatalf("Expected failure exit code, but got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "invalid value") {
		t.Fatalf("Expected invalid duration error, but got %q", stderr.String())
	}
}

func TestRunSearchRejectsZeroTimeout(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{"search", "--query", "test", "--timeout", "0"}, &stdout, &stderr)
	if exitCode != ExitFailure {
		t.Fatalf("Expected failure exit code, but got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "expected a positive duration") {
		t.Fatalf("Expected positive duration error, but got %q", stderr.String())
	}
}

func TestResolveSearchQueryAcceptsPositionalWords(t *testing.T) {
	query, err := resolveSearchQuery("", []string{"find", "bbolt", "documentation"})
	if err != nil {
		t.Fatalf("resolveSearchQuery returned error: %v", err)
	}
	if query != "find bbolt documentation" {
		t.Fatalf("Expected joined positional query, but got %q", query)
	}
}

func TestResolveSearchQueryRejectsBothForms(t *testing.T) {
	_, err := resolveSearchQuery("flag query", []string{"positional query"})
	if err == nil || !strings.Contains(err.Error(), "not both") {
		t.Fatalf("Expected conflicting query forms error, but got %v", err)
	}
}
