package imageProcessing

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func ApplyChromaTransparency(imageBytes []byte) ([]byte, error) {
	decodedImage, imageFormat, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode generated image for transparency post-processing: %w", err)
	}
	if imageFormat != "png" {
		return nil, fmt.Errorf("transparent post-processing expected PNG output, but got %q", imageFormat)
	}

	imageBounds := decodedImage.Bounds()
	transparentImage := image.NewNRGBA(imageBounds)
	draw.Draw(transparentImage, imageBounds, decodedImage, imageBounds.Min, draw.Src)
	makeEdgeChromaTransparent(transparentImage)

	output := bytes.NewBuffer(nil)
	if err := png.Encode(output, transparentImage); err != nil {
		return nil, fmt.Errorf("failed to encode transparent PNG output: %w", err)
	}

	return output.Bytes(), nil
}

func makeEdgeChromaTransparent(targetImage *image.NRGBA) {
	bounds := targetImage.Bounds()
	queue := make([]image.Point, 0, bounds.Dx()*2+bounds.Dy()*2)
	visited := make(map[image.Point]bool, bounds.Dx()*bounds.Dy()/2)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		enqueueChromaPoint(targetImage, image.Point{X: x, Y: bounds.Min.Y}, visited, &queue)
		enqueueChromaPoint(targetImage, image.Point{X: x, Y: bounds.Max.Y - 1}, visited, &queue)
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		enqueueChromaPoint(targetImage, image.Point{X: bounds.Min.X, Y: y}, visited, &queue)
		enqueueChromaPoint(targetImage, image.Point{X: bounds.Max.X - 1, Y: y}, visited, &queue)
	}

	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]
		setPointTransparent(targetImage, point)

		enqueueChromaPoint(targetImage, image.Point{X: point.X + 1, Y: point.Y}, visited, &queue)
		enqueueChromaPoint(targetImage, image.Point{X: point.X - 1, Y: point.Y}, visited, &queue)
		enqueueChromaPoint(targetImage, image.Point{X: point.X, Y: point.Y + 1}, visited, &queue)
		enqueueChromaPoint(targetImage, image.Point{X: point.X, Y: point.Y - 1}, visited, &queue)
	}
}

func enqueueChromaPoint(targetImage *image.NRGBA, point image.Point, visited map[image.Point]bool, queue *[]image.Point) {
	if !point.In(targetImage.Bounds()) || visited[point] {
		return
	}
	visited[point] = true
	if !isChromaColor(targetImage.NRGBAAt(point.X, point.Y)) {
		return
	}
	*queue = append(*queue, point)
}

func isChromaColor(pixel color.NRGBA) bool {
	redDistance := absInt(int(pixel.R) - ChromaRed)
	greenDistance := absInt(int(pixel.G) - ChromaGreen)
	blueDistance := absInt(int(pixel.B) - ChromaBlue)
	if redDistance+greenDistance+blueDistance <= ChromaDistanceTolerance {
		return true
	}

	return int(pixel.G) >= ChromaMinimumGreen &&
		int(pixel.G) > int(pixel.R)+ChromaDominanceTolerance &&
		int(pixel.G) > int(pixel.B)+ChromaDominanceTolerance
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
