package logger

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	logger *log.Logger
	env    string
}

type LogEntry struct {
	Level     string      `json:"level"`
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	File      string      `json:"file,omitempty"`
	Line      int         `json:"line,omitempty"`
	Error     string      `json:"error,omitempty"`
}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(os.Getenv("ENVIRONMENT"))
}

func NewLogger(env string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
		env:    env,
	}
}

func (l *Logger) log(level string, message string, data interface{}, err error) {
	entry := LogEntry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
		Data:      data,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	if l.env == "development" {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.File = file
			entry.Line = line
		}
	}

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf("Error marshaling log entry: %v", err)
		return
	}

	l.logger.Println(string(jsonEntry))
}

func Info(ctx context.Context, message string, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	defaultLogger.log("INFO", message, logData, nil)
}

func Error(ctx context.Context, message string, err error, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		logData = data[0]
	}
	defaultLogger.log("ERROR", message, logData, err)
}

func Debug(ctx context.Context, message string, data ...interface{}) {
	if defaultLogger.env == "development" {
		var logData interface{}
		if len(data) > 0 {
			logData = data[0]
		}
		defaultLogger.log("DEBUG", message, logData, nil)
	}
}
