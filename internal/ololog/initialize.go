package ololog

import (
	"os"

	"github.com/rs/zerolog"
)

var lgr zerolog.Logger

func init() {
	// TODO
	// Думаю, тут должен создаваться новый экземпляр журнала, чтобы не дай бор там не пересечься с чем-нить ещё с зерологгерским.
	// все равно у нас же синглтон!
	//
	// Плюс должны быть настройки - вывод в файл, или что выводить в консоль

	//цветной вывод только если указана переменная окружения дебуг, если нет, то выводим тупо в файл
	lgr = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
