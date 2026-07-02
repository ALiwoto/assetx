package imageProcessing

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"strconv"
	"strings"
)

func ApplyChromaTransparency(imageBytes []byte) ([]byte, error) {
	chromaColor, err := ParseHexColor(ChromaHexColor)
	if err != nil {
		return nil, err
	}

	return RemoveColorTransparency(imageBytes, chromaColor, ChromaDistanceTolerance)
}

func RemoveColorTransparency(imageBytes []byte, targetColor color.NRGBA, tolerance int) ([]byte, error) {
	if tolerance < 0 {
		return nil, fmt.Errorf("color removal tolerance cannot be negative, got %d", tolerance)
	}

	decodedImage, imageFormat, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode generated image for transparency post-processing: %w", err)
	}
	if imageFormat == "" {
		return nil, fmt.Errorf("failed to detect image format for transparency post-processing")
	}

	imageBounds := decodedImage.Bounds()
	transparentImage := image.NewNRGBA(imageBounds)
	draw.Draw(transparentImage, imageBounds, decodedImage, imageBounds.Min, draw.Src)
	makeMatchingPixelsTransparent(transparentImage, targetColor, tolerance)

	output := bytes.NewBuffer(nil)
	if err := png.Encode(output, transparentImage); err != nil {
		return nil, fmt.Errorf("failed to encode transparent PNG output: %w", err)
	}

	return output.Bytes(), nil
}

func ParseHexColor(hexColor string) (color.NRGBA, error) {
	normalizedHexColor := strings.TrimSpace(hexColor)
	normalizedHexColor = strings.TrimPrefix(normalizedHexColor, "#")
	if len(normalizedHexColor) != 6 {
		return color.NRGBA{}, fmt.Errorf("invalid color %q: expected #RRGGBB", hexColor)
	}

	red, err := parseHexColorComponent(normalizedHexColor[0:2], hexColor)
	if err != nil {
		return color.NRGBA{}, err
	}
	green, err := parseHexColorComponent(normalizedHexColor[2:4], hexColor)
	if err != nil {
		return color.NRGBA{}, err
	}
	blue, err := parseHexColorComponent(normalizedHexColor[4:6], hexColor)
	if err != nil {
		return color.NRGBA{}, err
	}

	return color.NRGBA{R: red, G: green, B: blue, A: 255}, nil
}

func makeMatchingPixelsTransparent(targetImage *image.NRGBA, targetColor color.NRGBA, tolerance int) {
	bounds := targetImage.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isTargetColor(targetImage.NRGBAAt(x, y), targetColor, tolerance) {
				setPointTransparent(targetImage, image.Point{X: x, Y: y})
			}
		}
	}
}

func isTargetColor(pixel color.NRGBA, targetColor color.NRGBA, tolerance int) bool {
	if pixel.A == 0 {
		return false
	}

	redDistance := absInt(int(pixel.R) - int(targetColor.R))
	greenDistance := absInt(int(pixel.G) - int(targetColor.G))
	blueDistance := absInt(int(pixel.B) - int(targetColor.B))
	return redDistance+greenDistance+blueDistance <= tolerance
}

func setPointTransparent(targetImage *image.NRGBA, point image.Point) {
	pixel := targetImage.NRGBAAt(point.X, point.Y)
	pixel.A = 0
	targetImage.SetNRGBA(point.X, point.Y, pixel)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func parseHexColorComponent(hexComponent string, originalHexColor string) (uint8, error) {
	value, err := strconv.ParseUint(hexComponent, 16, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid color %q: expected #RRGGBB", originalHexColor)
	}

	return uint8(value), nil
}
