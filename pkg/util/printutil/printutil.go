package printutil

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	NoColor bool
)

func Info(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgCyan)
	fmt.Println(msg)
	color.Unset()
}

func Warning(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgYellow)
	fmt.Println(msg)
	color.Unset()
}

func Error(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgRed)
	fmt.Println(msg)
	color.Unset()
}

func Success(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgGreen)
	fmt.Println(msg)
	color.Unset()
}
