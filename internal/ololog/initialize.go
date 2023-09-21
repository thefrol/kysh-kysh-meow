package ololog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// TODO
	// Думаю, тут должен создаваться новый экземпляр журнала, чтобы не дай бор там не пересечься с чем-нить ещё с зерологгерским.
	// все равно у нас же синглтон!
	//
	// Плюс должны быть настройки - вывод в файл, или что выводить в консоль
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Warn().Msg("logger ok")
}
