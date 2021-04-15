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
	fmt.Print(msg)
	color.Unset()
}

func Warning(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgYellow)
	fmt.Print(msg)
	color.Unset()
}

func Danger(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgRed)
	fmt.Print(msg)
	color.Unset()
}

func Success(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgGreen)
	fmt.Print(msg)
	color.Unset()
}
