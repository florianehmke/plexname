package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Log levels..
const (
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
)

var (
	logger *Logger
	once   sync.Once
)

// Logger struct encapsulates the used logger.
type Logger struct {
	il *log.Logger
	wl *log.Logger
	el *log.Logger
}

// Initialize initializes the global logger.
func Initialize(logFilePath *string) {
	once.Do(func() {
		logger = NewLogger(logFilePath)
	})
}

// NewLogger creates a new logger.
func NewLogger(logFilePath *string) *Logger {
	iw := []io.Writer{os.Stdout}
	ww := []io.Writer{os.Stdout}
	ew := []io.Writer{os.Stderr}

	if logFilePath != nil {
		logFile, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("could not open log file: %v", err)
		}
		iw = append(iw, logFile)
		ww = append(ww, logFile)
		ew = append(ew, logFile)
	}

	return &Logger{
		il: log.New(createWriter(iw), INFO+" ", log.LstdFlags|log.LUTC),
		wl: log.New(createWriter(ww), WARNING+" ", log.LstdFlags|log.LUTC),
		el: log.New(createWriter(ew), ERROR+" ", log.LstdFlags|log.LUTC),
	}
}

// Info logs a info message.
func (l *Logger) Info(msg string) {
	l.il.Println(msg)
}

// Infof logs a formatted info message.
func (l *Logger) Infof(msgFormat string, args ...interface{}) {
	l.Info(fmt.Sprintf(msgFormat, args...))
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.wl.Println(msg)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(msgFormat string, args ...interface{}) {
	l.Warn(fmt.Sprintf(msgFormat, args...))
}

// Error logs a error message.
func (l *Logger) Error(msg string) {
	l.el.Println(msg)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(msgFormat string, args ...interface{}) {
	l.Error(fmt.Sprintf(msgFormat, args...))
}

// Fatal logs a fatal message.
func (l *Logger) Fatal(msg string) {
	l.el.Fatal(msg)
}

// Fatalf logs a formatted fatal message.
func (l *Logger) Fatalf(msgFormat string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(msgFormat, args...))
}

// Info logs a info message.
func Info(msg string) {
	getInstance().Info(msg)
}

// Infof logs a formatted info message.
func Infof(msgFormat string, args ...interface{}) {
	getInstance().Infof(msgFormat, args...)
}

// Warn logs a warning message.
func Warn(msg string) {
	getInstance().Warn(msg)
}

// Warnf logs a formatted warning message.
func Warnf(msgFormat string, args ...interface{}) {
	getInstance().Warnf(msgFormat, args...)
}

// Error logs a error message.
func Error(msg string) {
	getInstance().Error(msg)
}

// Errorf logs a formatted error message.
func Errorf(msgFormat string, args ...interface{}) {
	getInstance().Errorf(msgFormat, args...)
}

// Fatal logs a fatal message.
func Fatal(msg string) {
	getInstance().Fatal(msg)
}

// Fatalf logs a formatted fatal message.
func Fatalf(msgFormat string, args ...interface{}) {
	getInstance().Fatalf(msgFormat, args...)
}

func createWriter(writers []io.Writer) io.Writer {
	if len(writers) == 1 {
		return io.Writer(writers[0])
	}
	return io.MultiWriter(writers...)
}

func getInstance() *Logger {
	once.Do(func() {
		logger = NewLogger(nil)
	})
	return logger
}
