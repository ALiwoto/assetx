package cli

import (
	"fmt"
	"io"
	"strings"
)

func HandleHelpCommand(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		PrintRootHelp(stdout)
		return ExitSuccess
	}
	if len(args) > 1 {
		_, _ = fmt.Fprintf(stderr, "assetx help: unexpected arguments: %s\n", strings.Join(args[1:], " "))
		return ExitFailure
	}

	switch args[0] {
	case "convert-tgs":
		PrintConvertTGSHelp(stdout)
		return ExitSuccess
	case "convert-webp":
		PrintConvertWEBPHelp(stdout)
		return ExitSuccess
	case "image":
		PrintImageHelp(stdout)
		return ExitSuccess
	case "remove-bg":
		PrintRemoveBackgroundHelp(stdout)
		return ExitSuccess
	case "version":
		PrintVersionHelp(stdout)
		return ExitSuccess
	case "config":
		PrintConfigHelp(stdout)
		return ExitSuccess
	default:
		_, _ = fmt.Fprintf(stderr, "assetx help: unknown help topic %q\n", args[0])
		PrintRootHelp(stderr)
		return ExitFailure
	}
}
