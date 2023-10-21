// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	//Создать хранилище
	store = storage.New()
}

var defaultConfig = config{
	Addr: ":8080",
}

func main() {
	cfg := configure(defaultConfig)

	ololog.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, MeowRouter())
	if err != nil {
		ololog.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}
}
