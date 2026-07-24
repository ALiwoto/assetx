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
}
