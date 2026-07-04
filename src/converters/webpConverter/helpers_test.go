package webpConverter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertWEBPToPNGRejectsInvalidWebP(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "input.webp")
	outputPath := filepath.Join(t.TempDir(), "output.png")
	if err := os.WriteFile(inputPath, []byte("not a webp image"), 0644); err != nil {
		t.Fatalf("Failed to write invalid WebP input: %v", err)
	}

	err := ConvertWebpToPNG(&ConvertWEBPOptions{
		InputPath:  inputPath,
		OutputPath: outputPath,
	})
	if err == nil {
		t.Fatal("Expected invalid WebP decode error, but got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode input WebP") {
		t.Fatalf("Expected decode error, but got %q", err.Error())
	}
}

func TestNormalizeConvertWEBPOptionsRejectsNonWEBPInput(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "input.png")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test input: %v", err)
	}

	err := normalizeConvertWEBPOptions(&ConvertWEBPOptions{
		InputPath:  inputPath,
		OutputPath: "output.png",
	})
	if err == nil {
		t.Fatal("Expected non-WebP input error, but got nil")
	}
	if !strings.Contains(err.Error(), ".webp input path") {
		t.Fatalf("Expected WebP input error, but got %q", err.Error())
	}
}

func TestNormalizeConvertWEBPOptionsRejectsNonPNGOutput(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "input.webp")
	if err := os.WriteFile(inputPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test input: %v", err)
	}

	err := normalizeConvertWEBPOptions(&ConvertWEBPOptions{
		InputPath:  inputPath,
		OutputPath: "output.jpg",
	})
	if err == nil {
		t.Fatal("Expected non-PNG output error, but got nil")
	}
	if !strings.Contains(err.Error(), ".png output path") {
		t.Fatalf("Expected PNG output error, but got %q", err.Error())
	}
}
