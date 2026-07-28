package cli

import (
	"fmt"
	"strings"
)

func resolveSearchQuery(queryFlag string, positionalArgs []string) (string, error) {
	queryFlag = strings.TrimSpace(queryFlag)
	positionalQuery := strings.TrimSpace(strings.Join(positionalArgs, " "))
	if queryFlag != "" && positionalQuery != "" {
		return "", fmt.Errorf("use either --query or a positional search query, not both")
	}
	if positionalQuery != "" {
		return positionalQuery, nil
	}

	return queryFlag, nil
}
