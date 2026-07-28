package cli

import (
	"assetx/src/core/appRunner"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

func runSearchCommand(args []string, configPath string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 1 && isHelpCommand(args[0]) {
		PrintSearchHelp(stdout)
		return nil
	}

	var allowedDomains repeatedStringFlag
	searchFlags := flag.NewFlagSet("assetx search", flag.ContinueOnError)
	searchFlags.SetOutput(stderr)
	searchFlags.String("query", "", "question or research request")
	searchFlags.String("model", appRunner.DefaultSearchModel, "Responses API model")
	searchFlags.String("context", appRunner.DefaultSearchContextSize, "web search context size: low, medium, or high")
	searchFlags.Var(&allowedDomains, "domain", "allowed search hostname without scheme; repeat for multiple domains")
	progressInterval := searchFlags.Duration(
		"progress-interval",
		appRunner.DefaultSearchProgressInterval,
		"stderr progress interval; use 0 to disable progress",
	)
	timeout := searchFlags.Duration("timeout", appRunner.DefaultSearchTimeout, "maximum time to wait for the Responses API")
	searchFlags.Usage = func() {
		PrintSearchHelp(stderr)
	}

	if err := searchFlags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if searchFlags.NArg() != 0 {
		return fmt.Errorf(
			"unexpected positional arguments: %s; pass the search text with --query \"...\"",
			strings.Join(searchFlags.Args(), " "),
		)
	}
	if *timeout <= 0 {
		return fmt.Errorf("invalid --timeout %q: expected a positive duration", *timeout)
	}

	request := &appRunner.SearchRequest{
		AllowedDomains:   []string(allowedDomains),
		ConfigPath:       configPath,
		ProgressInterval: *progressInterval,
		Timeout:          *timeout,
	}
	searchFlags.VisitAll(func(flagValue *flag.Flag) {
		switch flagValue.Name {
		case "context":
			request.SearchContextSize = flagValue.Value.String()
		case "model":
			request.Model = flagValue.Value.String()
		case "query":
			request.Query = flagValue.Value.String()
		}
	})

	return appRunner.RunSearch(request, stdout, stderr)
}
