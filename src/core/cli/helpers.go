package cli

import (
	"assetx/src/converters/tgsConverter"
	"assetx/src/core/appRunner"
	"assetx/src/core/imageProcessing"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	filteredArgs, configPath, err := extractConfigFlag(args)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "assetx: %v\n", err)
		return ExitFailure
	}

	if len(filteredArgs) == 0 {
		PrintRootHelp(stdout)
		return ExitSuccess
	}
	if isHelpCommand(filteredArgs[0]) {
		return HandleHelpCommand(filteredArgs[1:], stdout, stderr)
	}

	switch filteredArgs[0] {
	case "convert-tgs":
		if err := runConvertTGSCommand(filteredArgs[1:], stdout, stderr); err != nil {
			_, _ = fmt.Fprintf(stderr, "assetx convert-tgs: %v\n", err)
			return ExitFailure
		}
		return ExitSuccess
	case "convert-webp":
		if err := runConvertWEBPCommand(filteredArgs[1:], stdout, stderr); err != nil {
			_, _ = fmt.Fprintf(stderr, "assetx convert-webp: %v\n", err)
			return ExitFailure
		}
		return ExitSuccess
	case "image":
		if err := runImageCommand(filteredArgs[1:], configPath, stdout, stderr); err != nil {
			_, _ = fmt.Fprintf(stderr, "assetx image: %v\n", err)
			return ExitFailure
		}
		return ExitSuccess
	case "remove-bg":
		if err := runRemoveBackgroundCommand(filteredArgs[1:], stdout, stderr); err != nil {
			_, _ = fmt.Fprintf(stderr, "assetx remove-bg: %v\n", err)
			return ExitFailure
		}
		return ExitSuccess
	case "version", "--version", "-v":
		if len(filteredArgs) > 1 {
			_, _ = fmt.Fprintf(stderr, "assetx version: unexpected arguments: %s\n", strings.Join(filteredArgs[1:], " "))
			return ExitFailure
		}
		HandleVersionCommand(stdout)
		return ExitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "assetx: unknown command %q\n", filteredArgs[0])
		PrintRootHelp(stderr)
		return ExitFailure
	}
}

func runConvertTGSCommand(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintConvertTGSHelp(stdout)
		return nil
	}

	convertTGSFlags := flag.NewFlagSet("assetx convert-tgs", flag.ContinueOnError)
	convertTGSFlags.SetOutput(stderr)
	convertTGSFlags.String("in", "", "input Telegram .tgs or WebM emoji/sticker path")
	convertTGSFlags.String("out", "", "output sprite PNG path")
	convertTGSFlags.String("ffmpeg", tgsConverter.DefaultFFMPEGExecutable, "ffmpeg executable path")
	convertTGSFlags.Usage = func() {
		PrintConvertTGSHelp(stderr)
	}

	if err := convertTGSFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if convertTGSFlags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(convertTGSFlags.Args(), " "))
	}

	request := &appRunner.ConvertTGSRequest{}
	convertTGSFlags.VisitAll(func(flagValue *flag.Flag) {
		switch flagValue.Name {
		case "ffmpeg":
			request.FFMPEGPath = flagValue.Value.String()
		case "in":
			request.InputPath = flagValue.Value.String()
		case "out":
			request.OutputPath = flagValue.Value.String()
		}
	})

	return appRunner.RunConvertTGS(request, stdout)
}

func runConvertWEBPCommand(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintConvertWEBPHelp(stdout)
		return nil
	}

	convertWEBPFlags := flag.NewFlagSet("assetx convert-webp", flag.ContinueOnError)
	convertWEBPFlags.SetOutput(stderr)
	convertWEBPFlags.String("in", "", "input WebP image path")
	convertWEBPFlags.String("out", "", "output PNG path")
	convertWEBPFlags.Usage = func() {
		PrintConvertWEBPHelp(stderr)
	}

	if err := convertWEBPFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if convertWEBPFlags.NArg() > 1 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(convertWEBPFlags.Args(), " "))
	}

	request := &appRunner.ConvertWEBPRequest{}
	convertWEBPFlags.VisitAll(func(flagValue *flag.Flag) {
		switch flagValue.Name {
		case "in":
			request.InputPath = flagValue.Value.String()
		case "out":
			request.OutputPath = flagValue.Value.String()
		}
	})

	if convertWEBPFlags.NArg() == 1 {
		if strings.TrimSpace(request.InputPath) != "" {
			return fmt.Errorf("use either --in or one positional input path, not both")
		}
		request.InputPath = convertWEBPFlags.Arg(0)
	}

	return appRunner.RunConvertWEBP(request, stdout)
}

func runRemoveBackgroundCommand(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintRemoveBackgroundHelp(stdout)
		return nil
	}

	removeBackgroundFlags := flag.NewFlagSet("assetx remove-bg", flag.ContinueOnError)
	removeBackgroundFlags.SetOutput(stderr)
	removeBackgroundFlags.String("in", "", "input image path")
	removeBackgroundFlags.String("out", "", "output PNG path")
	removeBackgroundFlags.String("color", imageProcessing.ChromaHexColor, "hex color to remove, in #RRGGBB format")
	removeBackgroundFlags.Int("tolerance", imageProcessing.ChromaDistanceTolerance, "sum RGB distance tolerance")
	removeBackgroundFlags.Usage = func() {
		PrintRemoveBackgroundHelp(stderr)
	}

	if err := removeBackgroundFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if removeBackgroundFlags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(removeBackgroundFlags.Args(), " "))
	}

	request := &appRunner.RemoveBackgroundRequest{}
	removeBackgroundFlags.VisitAll(func(flagValue *flag.Flag) {
		switch flagValue.Name {
		case "color":
			request.Color = flagValue.Value.String()
		case "in":
			request.InputPath = flagValue.Value.String()
		case "out":
			request.OutputPath = flagValue.Value.String()
		case "tolerance":
			tolerance, err := strconv.Atoi(flagValue.Value.String())
			if err == nil {
				request.Tolerance = tolerance
			}
		}
	})

	return appRunner.RunRemoveBackground(request, stdout)
}

func runImageCommand(args []string, configPath string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintImageHelp(stdout)
		return nil
	}

	var avoids repeatedStringFlag
	var exampleNotes repeatedStringFlag
	var examples repeatedStringFlag

	imageFlags := flag.NewFlagSet("assetx image", flag.ContinueOnError)
	imageFlags.SetOutput(stderr)
	imageFlags.String("model", appRunner.DefaultImageModel, "image model")
	imageFlags.String("background", appRunner.BackgroundAuto, "auto, opaque, or transparent")
	imageFlags.String("prompt", "", "image prompt")
	imageFlags.Var(&avoids, "avoid", "thing to avoid in the generated image; repeat for multiple avoid items")
	imageFlags.Var(&examples, "example", "input example image path; repeat for multiple examples")
	imageFlags.Var(&exampleNotes, "example-note", "description for the matching --example; repeat once per --example")
	imageFlags.String("quality", appRunner.DefaultImageQuality, "auto, low, medium, or high")
	imageFlags.String("size", appRunner.DefaultImageSize, "auto or WIDTHxHEIGHT")
	imageFlags.String("out", "", "output image path")
	imageFlags.Usage = func() {
		PrintImageHelp(stderr)
	}

	if err := imageFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if imageFlags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(imageFlags.Args(), " "))
	}

	request := &appRunner.ImageRequest{
		Avoids:       []string(avoids),
		ConfigPath:   configPath,
		ExampleNotes: []string(exampleNotes),
		Examples:     []string(examples),
	}
	imageFlags.VisitAll(func(flagValue *flag.Flag) {
		switch flagValue.Name {
		case "background":
			request.Background = flagValue.Value.String()
		case "model":
			request.Model = flagValue.Value.String()
		case "out":
			request.OutputPath = flagValue.Value.String()
		case "prompt":
			request.Prompt = flagValue.Value.String()
		case "quality":
			request.Quality = flagValue.Value.String()
		case "size":
			request.Size = flagValue.Value.String()
		}
	})

	return appRunner.RunImage(request, stdout)
}

func extractConfigFlag(args []string) ([]string, string, error) {
	filteredArgs := make([]string, 0, len(args))
	configPath := ""

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--config" {
			if index+1 >= len(args) {
				return nil, "", fmt.Errorf("missing value after --config")
			}
			configPath = args[index+1]
			index++
			continue
		}
		if strings.HasPrefix(arg, "--config=") {
			configPath = strings.TrimPrefix(arg, "--config=")
			if configPath == "" {
				return nil, "", fmt.Errorf("missing value after --config=")
			}
			continue
		}
		filteredArgs = append(filteredArgs, arg)
	}

	return filteredArgs, configPath, nil
}

func isHelpCommand(command string) bool {
	return command == "help" || command == "-h" || command == "--help"
}
