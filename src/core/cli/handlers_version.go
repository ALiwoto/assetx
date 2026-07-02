package cli

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

func HandleVersionCommand(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "assetx version %s %s/%s\n", CurrentAssetxVersion, runtime.GOOS, runtime.GOARCH)

	revision := "unknown"
	commitTime := "unknown"
	modified := false
	goVersion := runtime.Version()

	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				commitTime = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}
	}

	if len(revision) > 7 && revision != "unknown" {
		revision = revision[:7]
	}
	if modified {
		revision += " (dirty)"
	}

	_, _ = fmt.Fprintf(writer, "Commit: %s\n", revision)
	_, _ = fmt.Fprintf(writer, "Commit Date: %s\n", commitTime)
	_, _ = fmt.Fprintf(writer, "Go: %s\n", goVersion)
}
