package logging

import (
	"log/slog"
	"os"
)

// Logger provides structured logging capabilities
type Logger struct {
	*slog.Logger
}

// LogLevel represents available log levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// New creates a new structured logger
func New(level LogLevel) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// NewDefault creates a logger with INFO level
func NewDefault() *Logger {
	return New(LevelInfo)
}

// WithComponent returns a logger with component context
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{Logger: l.Logger.With("component", component)}
}

// WithOperation returns a logger with operation context
func (l *Logger) WithOperation(operation string) *Logger {
	return &Logger{Logger: l.Logger.With("operation", operation)}
}

// WithFile returns a logger with file context
func (l *Logger) WithFile(filename string) *Logger {
	return &Logger{Logger: l.Logger.With("file", filename)}
}

// LogOperation logs the start and completion of an operation
func (l *Logger) LogOperation(operation string, fn func() error) error {
	l.Info("Starting operation", "operation", operation)
	err := fn()
	if err != nil {
		l.Error("Operation failed", "operation", operation, "error", err)
		return err
	}
	l.Info("Operation completed", "operation", operation)
	return nil
}

// LogValidation logs validation attempts and results
func (l *Logger) LogValidation(item, field string, valid bool, err error) {
	if valid {
		l.Debug("Validation passed", "item", item, "field", field)
	} else {
		l.Warn("Validation failed", "item", item, "field", field, "error", err)
	}
}

// LogFileOperation logs file system operations
func (l *Logger) LogFileOperation(operation, path string, err error) {
	if err != nil {
		l.Error("File operation failed",
			"operation", operation,
			"path", path,
			"error", err)
	} else {
		l.Debug("File operation completed",
			"operation", operation,
			"path", path)
	}
}

// Global logger instance
var defaultLogger *Logger

func init() {
	defaultLogger = NewDefault()
}

// Default returns the default logger
func Default() *Logger {
	return defaultLogger
}

// SetDefault sets the default logger
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// Convenience functions using the default logger
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func WithComponent(component string) *Logger {
	return defaultLogger.WithComponent(component)
}

func WithOperation(operation string) *Logger {
	return defaultLogger.WithOperation(operation)
}

func WithFile(filename string) *Logger {
	return defaultLogger.WithFile(filename)
}
