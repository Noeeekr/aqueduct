package main

import (
	"log"
	"os"
)

func LogInfo(text string) {
	getLogInfo().Println(text)
}
func LogFatal(text string) {
	getLogFatal().Fatal(text)
}

func getLogInfo() *log.Logger {
	return log.New(os.Stdout, "[INFO] ", log.Lshortfile|log.Ltime)
}

func getLogFatal() *log.Logger {
	return log.New(os.Stderr, "[ERROR] ", log.Lshortfile|log.Ltime)
}
