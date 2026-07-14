package tgsConverter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertTGSToSpritePNG(ctx context.Context, options *ConvertTGSOptions) error {
	if ctx == nil {
		return fmt.Errorf("cannot convert TGS because context is nil")
	}
	if err := normalizeConvertTGSOptions(options); err != nil {
		return err
	}
	if err := rejectUnsupportedTGSMagic(options.InputPath); err != nil {
		return err
	}

	ffmpegPath, err := resolveFFMPEGPath(options.FFMPEGPath)
	if err != nil {
		return err
	}

	tempDirectory, err := os.MkdirTemp("", "assetx-tgs-frames-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary frame directory: %w", err)
	}
	defer os.RemoveAll(tempDirectory)

	framePattern := filepath.Join(tempDirectory, "frame_%08d.png")
	if err := extractTGSFrames(ctx, ffmpegPath, options.InputPath, framePattern); err != nil {
		return err
	}

	framePaths, err := filepath.Glob(filepath.Join(tempDirectory, "frame_*.png"))
	if err != nil {
		return fmt.Errorf("failed to list extracted PNG frames: %w", err)
	}
	if len(framePaths) == 0 {
		return fmt.Errorf("ffmpeg did not extract any frames from %q", options.InputPath)
	}

	if err := writeFrameSpritePNG(framePaths, options.OutputPath); err != nil {
		return err
	}

	return nil
}

func normalizeConvertTGSOptions(options *ConvertTGSOptions) error {
	if options == nil {
		return fmt.Errorf("cannot normalize TGS conversion options because options is nil")
	}

	options.FFMPEGPath = strings.TrimSpace(options.FFMPEGPath)
	if options.FFMPEGPath == "" {
		options.FFMPEGPath = DefaultFFMPEGExecutable
	}

	options.InputPath = strings.TrimSpace(options.InputPath)
	if options.InputPath == "" {
		return fmt.Errorf("missing input TGS path")
	}
	if _, err := os.Stat(options.InputPath); err != nil {
		return fmt.Errorf("failed to access input TGS %q: %w", options.InputPath, err)
	}

	options.OutputPath = strings.TrimSpace(options.OutputPath)
	if options.OutputPath == "" {
		return fmt.Errorf("missing output PNG path")
	}
	if strings.ToLower(filepath.Ext(options.OutputPath)) != ".png" {
		return fmt.Errorf("TGS conversion requires a .png output path because the sprite sheet is written as PNG")
	}

	return nil
}

func rejectUnsupportedTGSMagic(inputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input TGS %q: %w", inputPath, err)
	}
	defer inputFile.Close()

	header := make([]byte, 4)
	bytesRead, err := inputFile.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read input TGS header from %q: %w", inputPath, err)
	}
	if bytesRead < 2 {
		return fmt.Errorf("input TGS %q is too small to identify", inputPath)
	}
	if header[0] == 0x1F && header[1] == 0x8B {
		return fmt.Errorf(
			"input %q looks like gzip-compressed Lottie .tgs; convert-tgs currently supports Telegram WebM/VP9 emoji and stickers through ffmpeg, while Lottie .tgs needs a Lottie renderer",
			inputPath,
		)
	}

	return nil
}

func resolveFFMPEGPath(ffmpegPath string) (string, error) {
	resolvedPath, err := exec.LookPath(ffmpegPath)
	if err == nil {
		return resolvedPath, nil
	}

	return "", fmt.Errorf("convert-tgs requires ffmpeg, but %q was not found; install ffmpeg or pass --ffmpeg path/to/ffmpeg", ffmpegPath)
}

func extractTGSFrames(ctx context.Context, ffmpegPath string, inputPath string, framePattern string) error {
	command := exec.CommandContext(ctx, ffmpegPath, tgsFrameExtractionArguments(inputPath, framePattern)...)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		stderrText := strings.TrimSpace(stderr.String())
		if stderrText == "" {
			return fmt.Errorf("ffmpeg failed to extract PNG frames from %q: %w", inputPath, err)
		}
		return fmt.Errorf("ffmpeg failed to extract PNG frames from %q: %w: %s", inputPath, err, stderrText)
	}

	return nil
}

func tgsFrameExtractionArguments(inputPath string, framePattern string) []string {
	return []string{
		"-hide_banner",
		"-y",
		"-i",
		inputPath,
		"-pix_fmt",
		"rgba",
		framePattern,
	}
}

func writeFrameSpritePNG(framePaths []string, outputPath string) error {
	if len(framePaths) == 0 {
		return fmt.Errorf("expected at least one frame path, but got 0")
	}

	firstFrame, err := decodePNGFrame(framePaths[0])
	if err != nil {
		return err
	}
	frameBounds := firstFrame.Bounds()
	frameWidth := frameBounds.Dx()
	frameHeight := frameBounds.Dy()
	if frameWidth <= 0 || frameHeight <= 0 {
		return fmt.Errorf("expected positive frame dimensions, but got %dx%d", frameWidth, frameHeight)
	}

	columnCount := spriteColumnCount(len(framePaths))
	rowCount := (len(framePaths) + columnCount - 1) / columnCount
	spriteImage := image.NewNRGBA(image.Rect(0, 0, frameWidth*columnCount, frameHeight*rowCount))
	if err := drawSpriteFrame(spriteImage, firstFrame, frameWidth, frameHeight, 0); err != nil {
		return err
	}

	for frameIndex := 1; frameIndex < len(framePaths); frameIndex++ {
		frameImage, err := decodePNGFrame(framePaths[frameIndex])
		if err != nil {
			return err
		}
		if frameImage.Bounds().Dx() != frameWidth || frameImage.Bounds().Dy() != frameHeight {
			return fmt.Errorf(
				"expected extracted frame %q to be %dx%d, but got %dx%d",
				framePaths[frameIndex],
				frameWidth,
				frameHeight,
				frameImage.Bounds().Dx(),
				frameImage.Bounds().Dy(),
			)
		}
		if err := drawSpriteFrame(spriteImage, frameImage, frameWidth, frameHeight, frameIndex); err != nil {
			return err
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create sprite PNG %q: %w", outputPath, err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, spriteImage); err != nil {
		return fmt.Errorf("failed to encode sprite PNG %q: %w", outputPath, err)
	}

	return nil
}

func decodePNGFrame(framePath string) (image.Image, error) {
	frameFile, err := os.Open(framePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open extracted frame %q: %w", framePath, err)
	}
	defer frameFile.Close()

	frameImage, err := png.Decode(frameFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode extracted PNG frame %q: %w", framePath, err)
	}

	return frameImage, nil
}

func drawSpriteFrame(spriteImage draw.Image, frameImage image.Image, frameWidth int, frameHeight int, frameIndex int) error {
	columnCount := spriteImage.Bounds().Dx() / frameWidth
	if columnCount <= 0 {
		return fmt.Errorf("expected positive sprite column count, but got %d", columnCount)
	}

	columnIndex := frameIndex % columnCount
	rowIndex := frameIndex / columnCount
	destination := image.Rect(
		columnIndex*frameWidth,
		rowIndex*frameHeight,
		(columnIndex+1)*frameWidth,
		(rowIndex+1)*frameHeight,
	)
	draw.Draw(spriteImage, destination, frameImage, frameImage.Bounds().Min, draw.Src)

	return nil
}

func spriteColumnCount(frameCount int) int {
	if frameCount <= 0 {
		return 0
	}

	columnCount := 1
	for columnCount*columnCount < frameCount {
		columnCount++
	}

	return columnCount
}
