package logger

import (
	"log"
	"os"
)

type LogMode int

const (
	ModeDevelopment LogMode = iota
	ModeProduction
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Init(mode LogMode) {
	flags := log.LstdFlags | log.Lshortfile

	switch mode {
	case ModeProduction:
		Info = log.New(os.Stdout, "INFO: ", flags)
		Warn = log.New(os.Stdout, "WARN: ", flags)
		Error = log.New(os.Stderr, "ERROR: ", flags)
	default:
		Info = log.New(os.Stdout, "INFO: ", flags)
		Warn = log.New(os.Stdout, "WARN: ", flags)
		Error = log.New(os.Stderr, "ERROR: ", flags)
	}

	Info.Println("Logger initialized")
}
