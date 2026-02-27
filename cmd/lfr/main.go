package main

import (
	"github.com/lgdd/lfr-cli/internal/cmd"
	"github.com/lgdd/lfr-cli/pkg/util/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.Fatal(err.Error())
	}
}
