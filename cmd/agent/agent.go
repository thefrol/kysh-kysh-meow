package main

import (
	"compress/gzip"
	"path"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/report"
	"github.com/thefrol/kysh-kysh-meow/internal/scheduler"
)

func init() {
	// Добавим компрессию при отправке данных
	report.UseBeforeRequest(ApplyGZIP(20, gzip.BestCompression))
}

var defaultConfig = config{
	Addr:            "localhost:8080",
	ReportInterval:  10,
	PollingInterval: 2,
}

func main() {
	config := configure(defaultConfig)

	// Метрики собираются во временное хранилище s,
	// где они хранятся в сыром виде и готовы превратиться
	// в массив metrica.Metrica
	var s report.Stats

	// запуск планировщика
	c := scheduler.New()
	//собираем данные раз в pollingInterval
	c.AddJob(time.Duration(config.PollingInterval)*time.Second, func() {
		//Обновляем данные в хранилище, тут же увеличиваем счетчик
		s = report.Fetch()
	})
	// отправляем данные раз в repostInterval
	c.AddJob(time.Duration(config.ReportInterval)*time.Second, func() {
		//отправляем на сервер метрики из хранилища s
		err := report.Send(s.ToTransport(), Endpoint(config))
		if err != nil {
			ololog.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
			return
		}
	})

	// Запускаем планировщик, и он занимает поток
	c.Serve(200 * time.Millisecond)

}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(cfg config) string {
	return "http://" + path.Join(cfg.Addr, "update")
}
