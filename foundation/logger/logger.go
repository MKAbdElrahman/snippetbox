package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents different log levels.
type LogLevel string

// Predefined log levels.
const (
	InfoLevel  LogLevel = "INF"
	ErrorLevel LogLevel = "ERR"
)

var defaultLogColors = LogColors{
	InfoColor:      color.New(color.FgGreen),
	ErrorColor:     color.New(color.FgRed),
	TimestampColor: color.New(color.FgWhite),
	KeyColor:       color.New(color.FgHiBlack),
}

// LogColors contains the color configurations for the logger.
type LogColors struct {
	InfoColor      *color.Color
	ErrorColor     *color.Color
	TimestampColor *color.Color
	KeyColor       *color.Color
}

// Logger represents a simple logger with colored output.
type Logger struct {
	stdout         io.Writer
	stderr         io.Writer
	colors         LogColors
	withDate       map[LogLevel]bool
	withTime       map[LogLevel]bool
	withLineSource map[LogLevel]bool
}

// LoggerOption is a function type for configuring Logger.
type LoggerOption func(*Logger)

// WithColors configures custom colors for the logger.
func WithColors(colors LogColors) LoggerOption {
	return func(logger *Logger) {
		logger.colors = colors
	}
}

// WithLineSource configures whether to include the source code line number in log messages.
func WithLineSource(withLineSource bool, levels ...LogLevel) LoggerOption {
	return func(logger *Logger) {
		if len(levels) == 0 {
			logger.withLineSource = map[LogLevel]bool{
				"": withLineSource,
			}
		} else {
			for _, level := range levels {
				logger.withLineSource[level] = withLineSource
			}
		}
	}
}

// WithStdout configures a custom io.Writer for stdout.
func WithStdout(stdout io.Writer) LoggerOption {
	return func(logger *Logger) {
		logger.stdout = stdout
	}
}

// WithStderr configures a custom io.Writer for stderr.
func WithStderr(stderr io.Writer) LoggerOption {
	return func(logger *Logger) {
		logger.stderr = stderr
	}
}

// WithDate configures whether to include the date in log messages for specific log levels.
func WithDate(withDate bool, levels ...LogLevel) LoggerOption {
	return func(logger *Logger) {
		if len(levels) == 0 {
			logger.withDate = map[LogLevel]bool{
				"": withDate,
			}
		} else {
			for _, level := range levels {
				logger.withDate[level] = withDate
			}
		}
	}
}

// WithTime configures whether to include the time in log messages for specific log levels.
func WithTime(withTime bool, levels ...LogLevel) LoggerOption {
	return func(logger *Logger) {
		if len(levels) == 0 {
			logger.withTime = map[LogLevel]bool{
				"": withTime,
			}
		} else {
			for _, level := range levels {
				logger.withTime[level] = withTime
			}
		}
	}
}

// NewStdLogger creates a new instance of Logger with standard or custom configuration.
func NewLogger(options ...LoggerOption) *Logger {

	logger := &Logger{
		stdout:         os.Stdout,
		stderr:         os.Stderr,
		colors:         defaultLogColors,
		withDate:       make(map[LogLevel]bool),
		withTime:       make(map[LogLevel]bool),
		withLineSource: make(map[LogLevel]bool),
	}

	for _, option := range options {
		option(logger)
	}

	return logger
}

// Info logs an informational message.
func (logger *Logger) Info(msg string, fields ...map[string]interface{}) {
	fmt.Fprint(logger.stdout, logger.formatLogMessage(InfoLevel, msg, fields...))
}

// Error logs an error message.
func (logger *Logger) Error(msg string, fields ...map[string]interface{}) {
	fmt.Fprint(logger.stderr, logger.formatLogMessage(ErrorLevel, msg, fields...))
}

// formatLogMessage formats the log message with colors and additional fields.
func (logger *Logger) formatLogMessage(level LogLevel, msg string, fields ...map[string]interface{}) string {

	levelColor := logger.getLevelColor(level)

	timestamp := logger.getTimestamp(level)
	lineSource := logger.getLineSource(level)

	logMsg := fmt.Sprintf("[%s] %s: %s", logger.colors.TimestampColor.Sprint(timestamp), levelColor.Sprint(level), msg)

	if lineSource != "" {
		logMsg = fmt.Sprintf("[%s][%s] %s: %s", logger.colors.TimestampColor.Sprint(timestamp), lineSource, levelColor.Sprint(level), msg)
	}

	for _, field := range fields {
		for key, value := range field {
			ckey := logger.colors.KeyColor.Sprint(key)
			logMsg += fmt.Sprintf(" %s:%v", ckey, value)
		}
	}

	return logMsg + "\n"
}

// getTimestamp returns the formatted timestamp with date and/or time based on configuration.
func (logger *Logger) getTimestamp(level LogLevel) string {
	format := ""
	if withDate, ok := logger.withDate[level]; ok && withDate {
		format += "2006-01-02"
	}
	if withTime, ok := logger.withTime[level]; ok && withTime {
		if withDate := logger.withDate[level]; withDate {
			format += " "
		}
		format += "15:04:05"
	}

	return time.Now().Format(format)
}

// getLevelColor returns the color for a specific log level.
func (logger *Logger) getLevelColor(level LogLevel) *color.Color {
	switch level {
	case InfoLevel:
		return logger.colors.InfoColor
	case ErrorLevel:
		return logger.colors.ErrorColor
	default:
		return logger.colors.TimestampColor
	}
}

// getLineSource returns the formatted source code line number based on configuration.
func (logger *Logger) getLineSource(level LogLevel) string {
	if withLineSource, ok := logger.withLineSource[level]; ok && withLineSource {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			return fmt.Sprintf("%s:%d", filepath.Base(file), line)
		}
	}
	return ""
}
