package cli

import "strings"

func (flagValues *repeatedStringFlag) Set(value string) error {
	*flagValues = append(*flagValues, value)
	return nil
}

func (flagValues *repeatedStringFlag) String() string {
	return strings.Join(*flagValues, ",")
}
