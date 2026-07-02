package main

import (
	"assetx/src/core/cli"
	"os"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
