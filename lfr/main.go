package main

import (
	"fmt"
	"os"

	"github.com/lgdd/liferay-cli/lfr/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
