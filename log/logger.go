package log

import (
	"strings"
	"os"
	"log"
)

type Level uint8

const (
	Error Level = iota
	Warn
	Info
	Debug
)

type Logger struct {
 	Level Level
}

func Get() *Logger {

	var l = new(Logger)

	switch strings.ToLower(os.Getenv("REFINERY_LOGLEVEL")) {
	case "error":
		l.Level = Error
		break
	case "warn":
		l.Level = Warn
		break
	case "info":
		l.Level = Info
		break
	case "debug":
		l.Level = Debug
		break
	default:
		l.Level = Debug
	}

	return l
}

func printLog(message string) {
	log.Print(message)
}

func (l *Logger) Error(message string) {
	if (l.Level >= Error) {
		printLog(message)
	}
}

func (l*Logger) Warn(message string) {
	if (l.Level >= Warn) {
		printLog(message)
	}
}

func (l*Logger) Info(message string) {
	if (l.Level >= Info) {
		printLog(message)
	}
}

func (l*Logger) Debug(message string) {
	if (l.Level >= Debug) {
		printLog(message)
	}
}