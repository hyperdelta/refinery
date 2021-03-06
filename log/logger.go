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

func printLog(message interface{}) {
	log.Print(message)
}

func (l* Logger) Print(message interface{}) {
	printLog(message)
}

func (l *Logger) Error(message interface{}) {
	if (l.Level >= Error) {
		//msg, _ := json.Marshal(message)
		printLog(message)
	}
}

func (l*Logger) Warn(message interface{}) {
	if (l.Level >= Warn) {
		//msg, _ := json.Marshal(message)
		printLog(message)
	}
}

func (l*Logger) Info(message interface{}) {
	if (l.Level >= Info) {
		//msg, _ := json.Marshal(message)
		printLog(message)
	}
}

func (l*Logger) Debug(message interface{}) {
	if (l.Level >= Debug) {
		//msg, _ := json.Marshal(message)
		printLog(message)
	}
}