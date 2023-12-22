package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	Stdout         io.Writer
	Stderr         io.Writer
	InfoColor      *color.Color
	ErrorColor     *color.Color
	TimestampColor *color.Color
	KeyColor       *color.Color
}

func NewStdLogger() *Logger {
	return &Logger{
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		InfoColor:      color.New(color.FgGreen),
		ErrorColor:     color.New(color.FgRed),
		TimestampColor: color.New(color.FgYellow),
		KeyColor:       color.New(color.FgBlue),
	}
}

func (logger *Logger) Info(msg string, fields ...map[string]interface{}) {
	fmt.Fprint(logger.Stdout, logger.formatLogMessage(logger.InfoColor, "INF", msg, fields...))
}

func (logger *Logger) Error(msg string, fields ...map[string]interface{}) {
	fmt.Fprint(logger.Stderr, logger.formatLogMessage(logger.ErrorColor, "ERR", msg, fields...))
}

func (logger *Logger) formatLogMessage(color *color.Color, level, msg string, fields ...map[string]interface{}) string {
	timestamp := logger.TimestampColor.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, color.Sprint(level), msg)

	for _, field := range fields {
		for key, value := range field {
			ckey := logger.KeyColor.Sprint(key)
			logMsg += fmt.Sprintf(" %s:%v", ckey, value)
		}
	}

	return logMsg + "\n"
}
