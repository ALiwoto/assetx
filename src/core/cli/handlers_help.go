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
	case "image":
		PrintImageHelp(stdout)
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

func PrintRootHelp(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "assetx - generate game image assets with OpenAI image models.")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Usage:")
	_, _ = fmt.Fprintln(writer, "  assetx [--config path/to/config.json] <command> [options]")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Commands:")
	_, _ = fmt.Fprintln(writer, "  image      Generate or edit an image asset")
	_, _ = fmt.Fprintln(writer, "  help       Show help for assetx or a command")
	_, _ = fmt.Fprintln(writer, "  version    Print version, target platform, commit, and commit date")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Supported image models:")
	_, _ = fmt.Fprintln(writer, "  gpt-image-2, gpt-image-1.5")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Common help:")
	_, _ = fmt.Fprintln(writer, "  assetx help image")
	_, _ = fmt.Fprintln(writer, "  assetx help config")
	_, _ = fmt.Fprintln(writer, "  assetx image --help")
	_, _ = fmt.Fprintln(writer)
	PrintConfigSummary(writer)
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Example:")
	_, _ = fmt.Fprintln(writer, "  assetx image --model gpt-image-2 --background transparent --prompt \"create a battle win header\" --example example1.png --example example2.png --quality medium --size 1024x1024 --out assets/sprites/slime.png")
}

func PrintImageHelp(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "assetx image - generate or edit an image asset.")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Usage:")
	_, _ = fmt.Fprintln(writer, "  assetx [--config path/to/config.json] image --prompt \"...\" --out path/to/output.png [options]")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Required:")
	_, _ = fmt.Fprintln(writer, "  --prompt string        Text prompt for the asset")
	_, _ = fmt.Fprintln(writer, "  --out path             Output path ending in .png, .jpg, .jpeg, or .webp")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Options:")
	_, _ = fmt.Fprintln(writer, "  --model string         gpt-image-2 or gpt-image-1.5 (default: gpt-image-2)")
	_, _ = fmt.Fprintln(writer, "  --background string    auto, opaque, or transparent (default: auto)")
	_, _ = fmt.Fprintln(writer, "  --example path         Example/input image path; repeat for multiple examples")
	_, _ = fmt.Fprintln(writer, "  --quality string       auto, low, medium, or high (default: medium)")
	_, _ = fmt.Fprintln(writer, "  --size string          auto or WIDTHxHEIGHT (default: 1024x1024)")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Model and size rules:")
	_, _ = fmt.Fprintln(writer, "  gpt-image-2 accepts WIDTHxHEIGHT where both values are positive, divisible by 16, max 3840x2160, and aspect ratio is between 1:3 and 3:1.")
	_, _ = fmt.Fprintln(writer, "  gpt-image-1.5 accepts auto, 1024x1024, 1024x1536, or 1536x1024.")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Transparency:")
	_, _ = fmt.Fprintln(writer, "  gpt-image-1.5 can request transparent background directly.")
	_, _ = fmt.Fprintln(writer, "  gpt-image-2 does not support direct transparency; assetx requests a chroma background and post-processes edge-connected green pixels into PNG alpha.")
	_, _ = fmt.Fprintln(writer, "  gpt-image-2 transparent output requires --out ending in .png.")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Examples:")
	_, _ = fmt.Fprintln(writer, "  assetx image --prompt \"16x16 slime sprite, idle frame\" --out assets/sprites/slime.png")
	_, _ = fmt.Fprintln(writer, "  assetx image --model gpt-image-2 --background transparent --prompt \"battle win header\" --quality medium --size 1024x1024 --out assets/ui/battle_win.png")
	_, _ = fmt.Fprintln(writer, "  assetx image --model gpt-image-1.5 --background transparent --prompt \"match this style\" --example refs/style.png --out assets/icons/gem.png")
}

func PrintConfigHelp(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "assetx config - configuration file format.")
	_, _ = fmt.Fprintln(writer)
	PrintConfigSummary(writer)
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "JSON example:")
	_, _ = fmt.Fprintln(writer, "{")
	_, _ = fmt.Fprintln(writer, "  \"proxy_base_url\": \"\",")
	_, _ = fmt.Fprintln(writer, "  \"api_key\": \"b64::base64_encoded_key_here\"")
	_, _ = fmt.Fprintln(writer, "}")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Proxy example:")
	_, _ = fmt.Fprintln(writer, "  \"proxy_base_url\": \"https://main.purroxy.org/openai/v1\"")
}

func PrintVersionHelp(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "assetx version - print build identity.")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Usage:")
	_, _ = fmt.Fprintln(writer, "  assetx version")
	_, _ = fmt.Fprintln(writer, "  assetx --version")
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, "Output includes the assetx version, GOOS/GOARCH, VCS commit, dirty marker, commit date, and Go compiler version when available.")
}

func PrintConfigSummary(writer io.Writer) {
	_, _ = fmt.Fprintln(writer, "Config:")
	_, _ = fmt.Fprintln(writer, "  Default path: ~/.assetx/config.json")
	_, _ = fmt.Fprintln(writer, "  Override path: assetx --config some/other/path/config.json <command>")
	_, _ = fmt.Fprintln(writer, "  proxy_base_url empty string uses OpenAI directly at https://api.openai.com/v1.")
	_, _ = fmt.Fprintln(writer, "  proxy_base_url non-empty must be an OpenAI-compatible API base URL.")
	_, _ = fmt.Fprintln(writer, "  api_key may be raw or prefixed with b64:: followed by a base64-encoded key.")
}
