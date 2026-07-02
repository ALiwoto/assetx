package cli

import (
	"assetx/src/core/appRunner"
	"errors"
	"flag"
	"fmt"
	"io"
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
	case "image":
		if err := runImageCommand(filteredArgs[1:], configPath, stdout, stderr); err != nil {
			_, _ = fmt.Fprintf(stderr, "assetx image: %v\n", err)
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

func runImageCommand(args []string, configPath string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintImageHelp(stdout)
		return nil
	}

	var examples repeatedStringFlag

	imageFlags := flag.NewFlagSet("assetx image", flag.ContinueOnError)
	imageFlags.SetOutput(stderr)
	imageFlags.String("model", appRunner.DefaultImageModel, "OpenAI image model")
	imageFlags.String("background", appRunner.BackgroundAuto, "auto, opaque, or transparent")
	imageFlags.String("prompt", "", "image prompt")
	imageFlags.Var(&examples, "example", "input example image path; repeat for multiple examples")
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

	request := appRunner.ImageRequest{
		ConfigPath: configPath,
		Examples:   []string(examples),
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
