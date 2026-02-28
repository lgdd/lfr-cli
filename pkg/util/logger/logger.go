// Package logger provides styled terminal output using charmbracelet/lipgloss
// and charmbracelet/log. It exposes leveled log functions (Debug, Info, Warn,
// Error, Fatal) and styled print helpers for info, error, warn, and success
// messages. Color output is suppressed when conf.NoColor is true.
package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/lgdd/lfr-cli/internal/conf"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})
	boldStyle    = lipgloss.NewStyle().Bold(true)
	debugStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	infoStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("87"))
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("160"))
	fatalStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	warnStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("190"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("83"))
	noColorStyle = lipgloss.NewStyle().UnsetForeground().UnsetBackground()
)

// Debug logs a message at the debug level with optional key-value pairs.
func Debug(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Debug(msg, keyvals...)
}

// Info logs a message at the info level with optional key-value pairs.
func Info(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Info(msg, keyvals...)
}

// Warn logs a message at the warn level with optional key-value pairs.
func Warn(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Warn(msg, keyvals...)
}

// Error logs a message at the error level with optional key-value pairs.
func Error(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Error(msg, keyvals...)
}

// Fatal logs a message at the fatal level with optional key-value pairs, then exits the program.
func Fatal(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Fatal(msg, keyvals...)
}

// Debugf logs a formatted message at the debug level.
func Debugf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Debug(format, a...)
}

// Infof logs a formatted message at the info level.
func Infof(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Info(format, a...)
}

// Warnf logs a formatted message at the warn level.
func Warnf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Warn(format, a...)
}

// Errorf logs a formatted message at the error level.
func Errorf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Error(format, a...)
}

// Fatalf logs a formatted message at the fatal level, then exits the program.
func Fatalf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Fatal(format, a...)
}

// Print prints a plain message without any styling or level prefix.
func Print(msg string) {
	fmt.Print(msg)
}

// Println prints a plain message followed by a newline, without any styling or level prefix.
func Println(msg string) {
	fmt.Println(msg)
}

// Printf prints a formatted plain message without any styling or level prefix.
func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

// PrintBold prints a bold message.
func PrintBold(msg string) {
	fmt.Print(lipgloss.Style.Render(boldStyle, msg))
}

// PrintlnBold prints a bold message followed by a newline.
func PrintlnBold(msg string) {
	fmt.Println(lipgloss.Style.Render(boldStyle, msg))
}

// PrintfBold prints a bold formatted message.
func PrintfBold(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(lipgloss.Style.Render(boldStyle, msg))
}

// PrintInfo prints a message styled with the info color.
func PrintInfo(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(infoStyle, msg))
	}
}

// PrintfInfo prints a formatted message styled with the info color.
func PrintfInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(infoStyle, msg))
	}
}

// PrintlnInfo prints a formatted message styled with the info color, followed by a newline.
func PrintlnInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(infoStyle, msg))
	}
}

// PrintError prints a message styled with the error color.
func PrintError(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(errorStyle, msg))
	}
}

// PrintfError prints a formatted message styled with the error color.
func PrintfError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(errorStyle, msg))
	}
}

// PrintlnError prints a message styled with the error color, followed by a newline.
func PrintlnError(msg string) {
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(errorStyle, msg))
	}
}

// PrintWarn prints a message styled with the warn color.
func PrintWarn(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(warnStyle, msg))
	}
}

// PrintfWarn prints a formatted message styled with the warn color.
func PrintfWarn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(warnStyle, msg))
	}
}

// PrintlnWarn prints a message styled with the warn color, followed by a newline.
func PrintlnWarn(msg string) {
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(warnStyle, msg))
	}
}

// PrintSuccess prints a message styled with the success color.
func PrintSuccess(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(successStyle, msg))
	}
}

// PrintfSuccess prints a formatted message styled with the success color.
func PrintfSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(successStyle, msg))
	}
}

// PrintlnSuccess prints a message styled with the success color, followed by a newline.
func PrintlnSuccess(msg string) {
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(successStyle, msg))
	}
}

func defaultStyles() *log.Styles {
	return &log.Styles{
		Timestamp: lipgloss.NewStyle(),
		Caller:    lipgloss.NewStyle().Faint(true),
		Prefix:    boldStyle.Faint(true),
		Message:   lipgloss.NewStyle(),
		Key:       lipgloss.NewStyle().Faint(true),
		Value:     lipgloss.NewStyle(),
		Separator: lipgloss.NewStyle().Faint(true),
		Levels: map[log.Level]lipgloss.Style{
			log.DebugLevel: debugStyle.
				SetString(strings.ToUpper(log.DebugLevel.String())),
			log.InfoLevel: infoStyle.
				SetString(strings.ToUpper(log.InfoLevel.String())),
			log.WarnLevel: warnStyle.
				SetString(strings.ToUpper(log.WarnLevel.String())),
			log.ErrorLevel: errorStyle.
				SetString(strings.ToUpper(log.ErrorLevel.String())),
			log.FatalLevel: fatalStyle.
				SetString(strings.ToUpper(log.FatalLevel.String())),
		},
		Keys:   map[string]lipgloss.Style{},
		Values: map[string]lipgloss.Style{},
	}
}

func noColorStyles() *log.Styles {
	return &log.Styles{
		Timestamp: noColorStyle,
		Caller:    noColorStyle,
		Prefix:    noColorStyle,
		Message:   noColorStyle,
		Key:       noColorStyle,
		Value:     noColorStyle,
		Separator: noColorStyle,
		Levels: map[log.Level]lipgloss.Style{
			log.DebugLevel: boldStyle.
				SetString(strings.ToUpper(log.DebugLevel.String())).
				UnsetForeground().
				UnsetBackground(),
			log.InfoLevel: boldStyle.
				SetString(strings.ToUpper(log.InfoLevel.String())).
				UnsetForeground().
				UnsetBackground(),
			log.WarnLevel: boldStyle.
				SetString(strings.ToUpper(log.WarnLevel.String())).
				UnsetForeground().
				UnsetBackground(),
			log.ErrorLevel: boldStyle.
				SetString(strings.ToUpper(log.ErrorLevel.String())).
				UnsetForeground().
				UnsetBackground(),
			log.FatalLevel: boldStyle.
				SetString(strings.ToUpper(log.FatalLevel.String())).
				UnsetForeground().
				UnsetBackground(),
		},
		Keys:   map[string]lipgloss.Style{},
		Values: map[string]lipgloss.Style{},
	}
}
