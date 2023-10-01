// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

var defaultConfig = config{
	Addr:                 ":8080",
	StoreIntervalSeconds: 300,
	FileStoragePath:      "/tmp/metrics-db.json",
	Restore:              true,
}

func main() {
	cfg := configure(defaultConfig)

	// создаем хранилище
	s, err := fileStorage(cfg)
	if err != nil {
		ololog.Error().Msgf("Не удалось сконфигурировать сервер, по причине: %v", err)
		return
	}
	store = s

	// Запускаем сервер в отдельной гоурутине
	ololog.Info().Msgf("^.^ Мяу, сервер запускается по адресу %v!", cfg.Addr)
	srv := http.Server{Addr: cfg.Addr, Handler: MeowRouter()}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		ololog.Error().Msgf("^0^ не могу запустить сервер: %v \n", err)
	}

}
