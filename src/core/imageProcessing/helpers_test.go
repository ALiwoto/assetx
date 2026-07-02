package imageProcessing

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemoveColorTransparencyRemovesEnclosedChroma(t *testing.T) {
	sourceImage := image.NewNRGBA(image.Rect(0, 0, 3, 3))
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			sourceImage.SetNRGBA(x, y, color.NRGBA{R: 80, G: 20, B: 120, A: 255})
		}
	}
	sourceImage.SetNRGBA(1, 1, color.NRGBA{R: ChromaRed, G: ChromaGreen, B: ChromaBlue, A: 255})

	input := bytes.NewBuffer(nil)
	if err := png.Encode(input, sourceImage); err != nil {
		t.Fatalf("Failed to encode source image: %v", err)
	}

	targetColor, err := ParseHexColor(ChromaHexColor)
	if err != nil {
		t.Fatalf("ParseHexColor returned error: %v", err)
	}
	outputBytes, err := RemoveColorTransparency(input.Bytes(), targetColor, ChromaDistanceTolerance)
	if err != nil {
		t.Fatalf("RemoveColorTransparency returned error: %v", err)
	}

	outputImage, err := png.Decode(bytes.NewReader(outputBytes))
	if err != nil {
		t.Fatalf("Failed to decode output image: %v", err)
	}

	centerColor := color.NRGBAModel.Convert(outputImage.At(1, 1)).(color.NRGBA)
	if centerColor.A != 0 {
		t.Fatalf("Expected enclosed chroma pixel to be transparent, but got alpha %d", centerColor.A)
	}

	cornerColor := color.NRGBAModel.Convert(outputImage.At(0, 0)).(color.NRGBA)
	if cornerColor.A != 255 {
		t.Fatalf("Expected non-chroma corner pixel to stay opaque, but got alpha %d", cornerColor.A)
	}
}

func TestWriteFrameSpritePNGComposesGrid(t *testing.T) {
	tempDirectory := t.TempDir()
	framePaths := []string{
		writeTestPNGFrame(t, tempDirectory, "frame_00000001.png", color.NRGBA{R: 255, A: 255}),
		writeTestPNGFrame(t, tempDirectory, "frame_00000002.png", color.NRGBA{G: 255, A: 255}),
		writeTestPNGFrame(t, tempDirectory, "frame_00000003.png", color.NRGBA{B: 255, A: 255}),
	}
	outputPath := filepath.Join(tempDirectory, "sprite.png")

	if err := writeFrameSpritePNG(framePaths, outputPath); err != nil {
		t.Fatalf("writeFrameSpritePNG returned error: %v", err)
	}

	outputFile, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open sprite output: %v", err)
	}
	defer outputFile.Close()

	spriteImage, err := png.Decode(outputFile)
	if err != nil {
		t.Fatalf("Failed to decode sprite output: %v", err)
	}
	if spriteImage.Bounds().Dx() != 4 || spriteImage.Bounds().Dy() != 2 {
		t.Fatalf("Expected sprite dimensions 4x2, but got %dx%d", spriteImage.Bounds().Dx(), spriteImage.Bounds().Dy())
	}

	expectedColors := []struct {
		point image.Point
		color color.NRGBA
	}{
		{point: image.Point{X: 0, Y: 0}, color: color.NRGBA{R: 255, A: 255}},
		{point: image.Point{X: 2, Y: 0}, color: color.NRGBA{G: 255, A: 255}},
		{point: image.Point{X: 0, Y: 1}, color: color.NRGBA{B: 255, A: 255}},
		{point: image.Point{X: 3, Y: 1}, color: color.NRGBA{}},
	}
	for _, expectedColor := range expectedColors {
		actualColor := color.NRGBAModel.Convert(spriteImage.At(expectedColor.point.X, expectedColor.point.Y)).(color.NRGBA)
		if actualColor != expectedColor.color {
			t.Fatalf("Expected color at %v to be %#v, but got %#v", expectedColor.point, expectedColor.color, actualColor)
		}
	}
}

func TestRejectUnsupportedTGSMagicRejectsGzipLottie(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "animated.tgs")
	if err := os.WriteFile(inputPath, []byte{0x1F, 0x8B, 0x08, 0x00}, 0644); err != nil {
		t.Fatalf("Failed to write test TGS: %v", err)
	}

	err := rejectUnsupportedTGSMagic(inputPath)
	if err == nil {
		t.Fatal("Expected gzip Lottie rejection, but got nil")
	}
	if !strings.Contains(err.Error(), "gzip-compressed Lottie .tgs") {
		t.Fatalf("Expected gzip Lottie error, but got %q", err.Error())
	}
}

func writeTestPNGFrame(t *testing.T, directory string, fileName string, frameColor color.NRGBA) string {
	t.Helper()

	frameImage := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	for x := 0; x < 2; x++ {
		frameImage.SetNRGBA(x, 0, frameColor)
	}

	framePath := filepath.Join(directory, fileName)
	frameFile, err := os.Create(framePath)
	if err != nil {
		t.Fatalf("Failed to create test PNG frame: %v", err)
	}
	defer frameFile.Close()

	if err := png.Encode(frameFile, frameImage); err != nil {
		t.Fatalf("Failed to encode test PNG frame: %v", err)
	}

	return framePath
}
