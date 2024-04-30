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

func Debug(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Debug(msg, keyvals...)
}

func Info(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Info(msg, keyvals...)
}

func Warn(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Warn(msg, keyvals...)
}

func Error(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Error(msg, keyvals...)
}

func Fatal(msg interface{}, keyvals ...interface{}) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Fatal(msg, keyvals...)
}

func Debugf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Debug(format, a...)
}

func Infof(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Info(format, a...)
}

func Warnf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Warn(format, a...)
}

func Errorf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Error(format, a...)
}

func Fatalf(format string, a ...any) {
	logger.SetStyles(defaultStyles())
	if conf.NoColor {
		logger.SetStyles(noColorStyles())
	}
	logger.Fatal(format, a...)
}

func Print(msg string) {
	fmt.Print(msg)
}

func Println(msg string) {
	fmt.Println(msg)
}

func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

func PrintBold(msg string) {
	fmt.Print(lipgloss.Style.Render(boldStyle, msg))
}

func PrintlnBold(msg string) {
	fmt.Println(lipgloss.Style.Render(boldStyle, msg))
}

func PrintfBold(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Print(lipgloss.Style.Render(boldStyle, msg))
}

func PrintInfo(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(infoStyle, msg))
	}
}

func PrintfInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(infoStyle, msg))
	}
}

func PrintlnInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(infoStyle, msg))
	}
}

func PrintError(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(errorStyle, msg))
	}
}

func PrintfError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(errorStyle, msg))
	}
}

func PrintlnError(msg string) {
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(errorStyle, msg))
	}
}

func PrintWarn(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(warnStyle, msg))
	}
}

func PrintfWarn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(warnStyle, msg))
	}
}

func PrintlnWarn(msg string) {
	if conf.NoColor {
		fmt.Println(msg)
	} else {
		fmt.Println(lipgloss.Style.Render(warnStyle, msg))
	}
}

func PrintSuccess(msg string) {
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(successStyle, msg))
	}
}

func PrintfSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if conf.NoColor {
		fmt.Print(msg)
	} else {
		fmt.Print(lipgloss.Style.Render(successStyle, msg))
	}
}

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
