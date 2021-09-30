package printutil

import (
	"fmt"

	"github.com/fatih/color"
)

// NoColor allows to disable colors for printed messages, default is false
var (
	NoColor bool
)

// Display info level message
func Info(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgCyan)
	fmt.Print(msg)
	color.Unset()
}

// Display warning level message
func Warning(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgYellow)
	fmt.Print(msg)
	color.Unset()
}

// Display critical level message
func Danger(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgRed)
	fmt.Print(msg)
	color.Unset()
}

// Display success message
func Success(msg string) {
	color.NoColor = NoColor
	color.Set(color.FgGreen)
	fmt.Print(msg)
	color.Unset()
}
