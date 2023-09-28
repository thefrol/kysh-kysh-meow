package main

import (
	"compress/gzip"
	"fmt"
	"path"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/report"
	"github.com/thefrol/kysh-kysh-meow/internal/scheduler"
	"github.com/thefrol/kysh-kysh-meow/internal/stats"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	store = storage.New()

	report.UseBeforeRequest(report.ApplyGZIP(1023, gzip.BestCompression))
}

var defaultConfig = config{
	Addr:            "localhost:8080",
	ReportInterval:  10,
	PollingInterval: 2,
}

func main() {
	config := configure(defaultConfig)

	// запуск планировщика
	c := scheduler.New()
	//собираем данные раз в pollingInterval
	c.AddJob(time.Duration(config.PollingInterval)*time.Second, func() {
		//Обновляем данные в хранилище, тут же увеличиваем счетчик
		stats.Fetch(store)
	})
	// отправляем данные раз в repostInterval
	c.AddJob(time.Duration(config.ReportInterval)*time.Second, func() {
		//отправляем на сервер
		err := report.Send(store.Metricas(), "http://"+path.Join(config.Addr, "update"))
		if err != nil {
			fmt.Println("Попытка отправить метрики завершилась с  ошибками:")
			fmt.Print(err)
			return
		}
		// Сбрасываем PollCount
		// #TODO: в таком случае нужно проверить, что счетчик реально отправился,
		//		а не просто, или нам пофигу?)
		stats.DropPollCount(store)
	})

	c.Serve(200 * time.Millisecond)

}
