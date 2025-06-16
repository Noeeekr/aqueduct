package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger(folder *os.File) *Logger {
	logger := log.New(folder, "", 0)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) SetOutput(folder *os.File) {
	l.logger.SetOutput(folder)
}

func (l *Logger) Info(message string) {
	_, filename, line, ok := runtime.Caller(1)
	if ok && program.Info.Environment == EnvironmentDevelopment {
		message = filename + ":" + strconv.Itoa(line) + " | " + message
	}

	message = "[INFO]" + message
	fmt.Println(message)
	l.logger.Println(message)
}

// Logs to stdout and the defined log folder and returns the formated message in form of an error
func (l *Logger) Error(message string) error {
	_, filename, line, ok := runtime.Caller(1)
	if ok && program.Info.Environment == EnvironmentDevelopment {
		message = filename + ":" + strconv.Itoa(line) + " | " + message
	}
	message = "[ERROR] " + message

	fmt.Println(message)
	l.logger.Println(message)

	return errors.New(message)
}
