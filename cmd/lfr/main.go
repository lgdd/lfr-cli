package main

import (
	"fmt"
	"os"

	"github.com/lgdd/lfr-cli/internal/cmd"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			logger.Fatal(err.Error())
		}
		logger.Fatal(err.Error())
	}
}
