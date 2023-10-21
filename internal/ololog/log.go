package ololog

import (
	"github.com/rs/zerolog"
)

func Log() *zerolog.Event {
	return lgr.Log()
}

func Info() *zerolog.Event {
	return lgr.Info()
}

func Debug() *zerolog.Event {
	return lgr.Debug()
}

func Warn() *zerolog.Event {
	return lgr.Warn()
}

func Error() *zerolog.Event {
	return lgr.Error()
}

func Fatal() *zerolog.Event {
	return lgr.Fatal()
}
