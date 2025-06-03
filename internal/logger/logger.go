package logger

import (
	"log"
	"os"
	"sync"
)

// Logger is a simple wrapper around Go's standard log package.
// It provides different logging levels.
type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	mu          sync.Mutex // For thread-safe logging
}

var defaultLogger *Logger
var once sync.Once

// NewLogger initializes and returns a new Logger instance.
func NewLogger() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return defaultLogger
}

// Info logs an informational message.
func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.infoLogger.Printf(format, v...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warnLogger.Printf(format, v...)
}

// Error logs an error message.
func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorLogger.Printf(format, v...)
}

// Fatal logs an error message and then exits the application.
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorLogger.Fatalf(format, v...)
}

// Global logger functions for convenience
func Info(format string, v ...interface{}) {
	NewLogger().Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	NewLogger().Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	NewLogger().Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	NewLogger().Fatal(format, v...)
}
