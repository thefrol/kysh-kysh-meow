package main

import (
	"fmt"
	"path"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/report"
	"github.com/thefrol/kysh-kysh-meow/lib/scheduler"
)

const updateRoute = "/updates"

func main() {
	config := mustConfigure(defaultConfig)

	// Метрики собираются во временное хранилище s,
	// где они хранятся в сыром виде и готовы превратиться
	// в массив metrica.Metrica
	var s report.Stats

	// запуск планировщика
	c := scheduler.New()
	//собираем данные раз в pollingInterval
	c.AddJob(time.Duration(config.PollingInterval)*time.Second, func() {
		//Обновляем данные в хранилище, тут же увеличиваем счетчик
		log.Debug().Msg("Считывание метрик")
		s = report.Fetch()
	})
	// отправляем данные раз в repostInterval
	c.AddJob(time.Duration(config.ReportInterval)*time.Second, func() {
		//отправляем на сервер метрики из хранилища s
		err := report.Send(s.ToTransport(), Endpoint(config))
		if err != nil {
			log.Error().Msgf("Попытка отправить метрики завершилась с  ошибками: %v", err)
			return
		}
	})

	// Запускаем планировщик, и он занимает поток
	c.Serve(200 * time.Millisecond)

}

// Endpoint формирует точку, куда агент будет посылать все запросы на основе своей текущей конфигурации
func Endpoint(cfg config) string {
	return fmt.Sprintf("%s%s", "http://", path.Join(cfg.Addr, updateRoute))
}
