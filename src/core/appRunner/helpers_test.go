package appRunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeImageRequestRequiresOneNotePerExample(t *testing.T) {
	examplePath := createTestExampleFile(t, "screenshot.png")

	_, _, err := normalizeImageRequest(&ImageRequest{
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
	_, _, err := normalizeImageRequest(&ImageRequest{
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

	request := &ImageRequest{
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
	}
	_, _, err := normalizeImageRequest(request)
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

func TestNormalizeImageRequestAppendsAvoidsToPrompt(t *testing.T) {
	request := &ImageRequest{
		Avoids: []string{
			"watermark",
			"photorealistic style",
		},
		OutputPath: "out.png",
		Prompt:     "make a pixel art sword",
	}
	_, _, err := normalizeImageRequest(request)
	if err != nil {
		t.Fatalf("normalizeImageRequest returned error: %v", err)
	}

	expectedSnippets := []string{
		"Avoid:",
		"- watermark",
		"- photorealistic style",
	}
	for _, expectedSnippet := range expectedSnippets {
		if !strings.Contains(request.Prompt, expectedSnippet) {
			t.Fatalf("Expected prompt to contain %q, but got:\n%s", expectedSnippet, request.Prompt)
		}
	}
}

func TestNormalizeImageRequestRejectsEmptyAvoid(t *testing.T) {
	_, _, err := normalizeImageRequest(&ImageRequest{
		Avoids:     []string{" "},
		OutputPath: "out.png",
		Prompt:     "make an icon",
	})
	if err == nil {
		t.Fatal("Expected empty avoid error, but got nil")
	}
	if !strings.Contains(err.Error(), "--avoid values cannot be empty") {
		t.Fatalf("Expected empty avoid error, but got %q", err.Error())
	}
}

func TestNormalizeImageRequestRejectsNilRequest(t *testing.T) {
	_, _, err := normalizeImageRequest(nil)
	if err == nil {
		t.Fatal("Expected nil request error, but got nil")
	}
	if !strings.Contains(err.Error(), "request is nil") {
		t.Fatalf("Expected nil request error, but got %q", err.Error())
	}
}

func TestNormalizeConvertTGSRequestRejectsNonPNGOutput(t *testing.T) {
	inputPath := createTestExampleFile(t, "custom_emoji_1.tgs")

	err := normalizeConvertTGSRequest(&ConvertTGSRequest{
		InputPath:  inputPath,
		OutputPath: "out.jpg",
	})
	if err == nil {
		t.Fatal("Expected non-PNG output error, but got nil")
	}
	if !strings.Contains(err.Error(), ".png --out path") {
		t.Fatalf("Expected PNG output error, but got %q", err.Error())
	}
}

func TestNormalizeConvertTGSRequestDefaultsFFMPEGPath(t *testing.T) {
	inputPath := createTestExampleFile(t, "custom_emoji_1.tgs")
	request := &ConvertTGSRequest{
		InputPath:  inputPath,
		OutputPath: "out.png",
	}

	if err := normalizeConvertTGSRequest(request); err != nil {
		t.Fatalf("normalizeConvertTGSRequest returned error: %v", err)
	}
	if request.FFMPEGPath != "ffmpeg" {
		t.Fatalf("Expected default ffmpeg path, but got %q", request.FFMPEGPath)
	}
}

func TestNormalizeConvertWEBPRequestDefaultsOutputPath(t *testing.T) {
	inputPath := createTestExampleFile(t, "example.webp")
	request := &ConvertWEBPRequest{
		InputPath: inputPath,
	}

	if err := normalizeConvertWEBPRequest(request); err != nil {
		t.Fatalf("normalizeConvertWEBPRequest returned error: %v", err)
	}
	if request.OutputPath != inputPath+".png" {
		t.Fatalf("Expected default output path %q, but got %q", inputPath+".png", request.OutputPath)
	}
}

func TestNormalizeConvertWEBPRequestRejectsNonWEBPInput(t *testing.T) {
	inputPath := createTestExampleFile(t, "example.png")

	err := normalizeConvertWEBPRequest(&ConvertWEBPRequest{
		InputPath: inputPath,
	})
	if err == nil {
		t.Fatal("Expected non-WebP input error, but got nil")
	}
	if !strings.Contains(err.Error(), ".webp input path") {
		t.Fatalf("Expected WebP input error, but got %q", err.Error())
	}
}

func TestNormalizeConvertWEBPRequestRejectsNonPNGOutput(t *testing.T) {
	inputPath := createTestExampleFile(t, "example.webp")

	err := normalizeConvertWEBPRequest(&ConvertWEBPRequest{
		InputPath:  inputPath,
		OutputPath: "out.jpg",
	})
	if err == nil {
		t.Fatal("Expected non-PNG output error, but got nil")
	}
	if !strings.Contains(err.Error(), ".png --out path") {
		t.Fatalf("Expected PNG output error, but got %q", err.Error())
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
