package logger

import (
	"cloud.google.com/go/logging"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

// Debug is debug message (severity 100)
func (logger logger) Debug(args ...interface{}) {
	logger.log(int(logging.Debug), args...)
}

func (logger logger) Debugf(formta string, args ...interface{}) {
	logger.log(int(logging.Debug), args...)
}

// Info is information message (severity 200)
func (logger logger) Info(args ...interface{}) {
	logger.log(int(logging.Info), args...)
}

// Warn is warning message (severity 400)
func (logger logger) Warn(args ...interface{}) {
	logger.log(int(logging.Warning), args...)
}

// Error is error message (severity 500)
func (logger logger) Error(args ...interface{}) {
	logger.log(int(logging.Error), args...)
}

func (logger logger) log(severity int, args ...interface{}) {
	entry := map[string]interface{}{
		"time":     time.Now().Format(time.RFC3339Nano),
		"severity": severity,
		"message":  args,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
