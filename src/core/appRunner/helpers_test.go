package appRunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeImageRequestRequiresOneNotePerExample(t *testing.T) {
	examplePath := createTestExampleFile(t, "screenshot.png")

	_, _, _, err := normalizeImageRequest(ImageRequest{
		Examples:   []string{examplePath},
		OutputPath: "out.png",
		Prompt:     "make an icon",
	})
	if err == nil {
		t.Fatal("Expected missing example note error, but got nil")
	}
	if !strings.Contains(err.Error(), "one --example-note per --example") {
		t.Fatalf("Expected example note count error, but got %q", err.Error())
	}
}

func TestNormalizeImageRequestRejectsNotesWithoutExamples(t *testing.T) {
	_, _, _, err := normalizeImageRequest(ImageRequest{
		ExampleNotes: []string{"screenshot of my game"},
		OutputPath:   "out.png",
		Prompt:       "make an icon",
	})
	if err == nil {
		t.Fatal("Expected note without example error, but got nil")
	}
	if !strings.Contains(err.Error(), "requires a matching --example") {
		t.Fatalf("Expected note without example error, but got %q", err.Error())
	}
}

func TestNormalizeImageRequestAppendsExampleNotesToPrompt(t *testing.T) {
	screenshotPath := createTestExampleFile(t, "screenshot1.png")
	assetPath := createTestExampleFile(t, "sword.png")

	request, _, _, err := normalizeImageRequest(ImageRequest{
		ExampleNotes: []string{
			"screenshot of my game",
			"existing asset style reference",
		},
		Examples: []string{
			screenshotPath,
			assetPath,
		},
		OutputPath: "out.png",
		Prompt:     "make an axe",
	})
	if err != nil {
		t.Fatalf("normalizeImageRequest returned error: %v", err)
	}

	expectedSnippets := []string{
		"Provided image reference notes:",
		"- Reference file 1 (screenshot1.png): screenshot of my game",
		"- Reference file 2 (sword.png): existing asset style reference",
		"Use these notes when interpreting the provided image references.",
	}
	for _, expectedSnippet := range expectedSnippets {
		if !strings.Contains(request.Prompt, expectedSnippet) {
			t.Fatalf("Expected prompt to contain %q, but got:\n%s", expectedSnippet, request.Prompt)
		}
	}
}

func createTestExampleFile(t *testing.T, fileName string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), fileName)
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test example file %q: %v", filePath, err)
	}

	return filePath
}
