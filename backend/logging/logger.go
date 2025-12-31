package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	// Global logger to write both stdout and file
	Logger *log.Logger
	// Keep reference to log file so it stays open
	logFile *os.File
)

/*
Function to initialize the logger with both stdout and file logging
*/
func InitLogger(logFilePath string) error {
	// Ensure the directory exists (extract from the full path)
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Open log file and keep it in package variable
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	Logger = log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	log.Println("Logger initialized, writing to", logFilePath)
	return nil
}

func Close() {
	if logFile != nil {
		log.Println("Closing log file")
		logFile.Sync() // Flush any buffered data
		logFile.Close()
	}
}
