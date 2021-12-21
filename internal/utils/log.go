package utils

import (
	"log"
	"os"
)

type Logger struct {
	InfoLog log.Logger
	ErrorLog log.Logger
}

func NewLogger() *Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime | log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ltime | log.Lshortfile)

	return &Logger {
		InfoLog: *infoLog,
		ErrorLog: *errorLog,
	}
}