package imageProcessing

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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
