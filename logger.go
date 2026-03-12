package main

import (
	"log"
	"os"
)

func LogToFile() *os.File {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("cant create log file")
	}

	return logFile
}
