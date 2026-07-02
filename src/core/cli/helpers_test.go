package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpCommandIncludesOperationalDetails(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	exitCode := Run([]string{"help"}, stdout, stderr)
	if exitCode != ExitSuccess {
		t.Fatalf("Expected exit code %d, but got %d with stderr %q", ExitSuccess, exitCode, stderr.String())
	}

	helpText := stdout.String()
	requiredSnippets := []string{
		"image-generation providers",
		"Currently implemented provider:",
		"assetx help image",
		"assetx help config",
		"gpt-image-2",
		"gpt-image-1.5",
		"b64::",
		"proxy_base_url",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(helpText, snippet) {
			t.Fatalf("Expected help text to contain %q, but got:\n%s", snippet, helpText)
		}
	}
	if strings.Contains(helpText, "with OpenAI image models") {
		t.Fatalf("Expected provider-neutral root help, but got:\n%s", helpText)
	}
}

func TestImageHelpCommandExitsSuccessfully(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	exitCode := Run([]string{"image", "--help"}, stdout, stderr)
	if exitCode != ExitSuccess {
		t.Fatalf("Expected exit code %d, but got %d with stderr %q", ExitSuccess, exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Transparency:") {
		t.Fatalf("Expected image help to describe transparency, but got:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Repeat --example") {
		t.Fatalf("Expected image help to describe multiple examples, but got:\n%s", stdout.String())
	}
}

func TestVersionCommandIncludesCommitFields(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	exitCode := Run([]string{"version"}, stdout, stderr)
	if exitCode != ExitSuccess {
		t.Fatalf("Expected exit code %d, but got %d with stderr %q", ExitSuccess, exitCode, stderr.String())
	}

	versionText := stdout.String()
	requiredSnippets := []string{
		"assetx version",
		"Commit:",
		"Commit Date:",
	}
	for _, snippet := range requiredSnippets {
		if !strings.Contains(versionText, snippet) {
			t.Fatalf("Expected version output to contain %q, but got:\n%s", snippet, versionText)
		}
	}
}
