package util

import (
	"fmt"
	"os"
)

const (
	InfoColor    = "\033[0;96m%s\033[0m"
	WarningColor = "\033[0;93m%s\033[0m"
	ErrorColor   = "\033[0;91m%s\033[0m"
	SuccessColor = "\033[0;92m%s\033[0m"
)

func PrintInfo(message string) {
	fmt.Printf(InfoColor, message)
	fmt.Println("")
}

func PrintWarning(message string) {
	fmt.Printf(WarningColor, message)
	fmt.Println("")
}

func PrintError(message string) {
	_, err := fmt.Fprintf(os.Stderr, ErrorColor, message)
	if err != nil {
		fmt.Printf(ErrorColor, message)
	}
	fmt.Println("")
}

func PrintSuccess(message string) {
	fmt.Printf(SuccessColor, message)
	fmt.Println("")
}
