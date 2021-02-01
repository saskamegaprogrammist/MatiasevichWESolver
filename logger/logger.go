package logger

import (
	"github.com/google/logger"
	"os"
)

var filename = "output.log"
var verbose = false
var logFile *os.File

func LoggerSetup() {
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("failed to open log file: %v", err)
	}
	logger.Init("MatiasevichLog", verbose, false, logFile)
}

func loggerClose() {
	logFile.Close()
	logger.Close()
}
