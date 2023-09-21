package ololog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Log() *zerolog.Event {
	return log.Log()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}
