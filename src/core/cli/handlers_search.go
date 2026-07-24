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
		return fmt.Errorf("unexpected positional arguments: %s", strings.Join(searchFlags.Args(), " "))
	}

	request := &appRunner.SearchRequest{
		AllowedDomains: []string(allowedDomains),
		ConfigPath:     configPath,
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

	return appRunner.RunSearch(request, stdout)
}
