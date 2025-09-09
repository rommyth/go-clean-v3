package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	TimeFormat = "2006-01-02 15:04:05"
)

func Init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: TimeFormat}
	logger := zerolog.New(output).With().Timestamp().Logger()
	log.Logger = logger
}

func Info(message string, fields map[string]interface{}) {
	evt := log.Info()
	for k,v := range fields {
		evt = evt.Interface(k, v)
	}
	evt.Msg(message)
}

// Error logs an error message with optional fields
func Error(message string, fields map[string]interface{}) {
    evt := log.Error()
    for k, v := range fields {
        evt = evt.Interface(k, v)
    }
    evt.Msg(message)
}

// Fatal logs a fatal message and exits
func Fatal(message string, fields map[string]interface{}) {
    evt := log.Fatal()
    for k, v := range fields {
        evt = evt.Interface(k, v)
    }
    evt.Msg(message)
}

// Debug logs a debug message (enable with LOG_LEVEL=debug)
func Debug(message string, fields map[string]interface{}) {
    evt := log.Debug()
    for k, v := range fields {
        evt = evt.Interface(k, v)
    }
    evt.Msg(message)
}