package logger

import (
	"fmt"
	"go-rest-template/pkg/config"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var logger *zerolog.Logger
var once sync.Once

func Initialize() {
	once.Do(func() {
		setupLogger()
	})
}

func setupLogger() {
	env := config.APP().GO_ENV
	debug := config.APP().DEBUG

	// Configure global settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return fmt.Sprintf("%s:%d", short, line)
	}

	var l zerolog.Logger

	if env == "development" || debug {
		// Development or debug
		l = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().CallerWithSkipFrameCount(2).Logger()

		l = l.Level(zerolog.DebugLevel)
	} else {
		// Production
		l = zerolog.New(os.Stdout).With().Timestamp().Logger()
		l = l.Level(zerolog.InfoLevel)
	}

	logger = &l
}

// Simple logging functions

func Debug(v ...interface{}) {
	Initialize()
	logger.Debug().CallerSkipFrame(1).Msgf("%v", v)
}

func Info(v ...interface{}) {
	Initialize()
	logger.Info().CallerSkipFrame(1).Msgf("%v", v)
}

func Error(v ...interface{}) {
	Initialize()
	logger.Error().CallerSkipFrame(1).Msgf("%v", v)
}

func Panic(v ...interface{}) {
	Initialize()
	logger.Panic().CallerSkipFrame(1).Msgf("%v", v)
}

// Formatted logging functions

func DebugF(format string, v ...interface{}) {
	Initialize()
	logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

func InfoF(format string, v ...interface{}) {
	Initialize()
	logger.Info().CallerSkipFrame(1).Msgf(format, v...)
}

func ErrorF(format string, v ...interface{}) {
	Initialize()
	logger.Error().CallerSkipFrame(1).Msgf(format, v...)
}

func PanicF(format string, v ...interface{}) {
	Initialize()
	logger.Panic().CallerSkipFrame(1).Msgf(format, v...)
}
