package log

import (
	log "github.com/sirupsen/logrus"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"

	UnrecognizedLogLevel = "Unrecognized log level"
)

func Debug(message string, keysAndValues ...interface{}) {
	Send(DebugLevel, message, keysAndValues...)
}

func Info(message string, keysAndValues ...interface{}) {
	Send(InfoLevel, message, keysAndValues...)
}

func Warn(message string, keysAndValues ...interface{}) {
	Send(WarnLevel, message, keysAndValues...)
}

func Error(message string, keysAndValues ...interface{}) {
	Send(ErrorLevel, message, keysAndValues...)
}

func Send(level LogLevel, message string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		send(level, message)
	}
	if len(keysAndValues)%2 != 0 {
		log.Error("Debug function must be given key-value pairs, the number of parameters must be even", "message", message)
		return
	}

	fields := make(log.Fields)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			log.Error("Key must be a string")
			return
		}
		fields[key] = keysAndValues[i+1]
	}

	sendWithFields(level, message, fields)
}

func send(level LogLevel, message string) {
	switch level {
	case DebugLevel:
		log.Debug(message)
	case InfoLevel:
		log.Info(message)
	case WarnLevel:
		log.Warn(message)
	case ErrorLevel:
		log.Error(message)
	default:
		log.Info(UnrecognizedLogLevel)
	}
}

func sendWithFields(level LogLevel, message string, fields log.Fields) {
	switch level {
	case DebugLevel:
		log.WithFields(fields).Debug(message)
	case InfoLevel:
		log.WithFields(fields).Info(message)
	case WarnLevel:
		log.WithFields(fields).Warn(message)
	case ErrorLevel:
		log.WithFields(fields).Error(message)
	default:
		log.WithFields(fields).Info(UnrecognizedLogLevel)
	}
}
