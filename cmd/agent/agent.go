package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/scheduler"
	"github.com/thefrol/kysh-kysh-meow/internal/stats"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	store = storage.New()
}

func main() {
	configure()

	// запуск планировщика
	c := scheduler.New()
	//собираем данные раз в pollingInterval
	c.AddJob(time.Duration(*pollIntervalSeconds)*time.Second, func() {
		//Обновляем данные в хранилище
		stats.FetchMemStats(store)
		// Увеличиваем PollCount
		incrementCounter(store, metricPollCount)
	})
	// отправляем данные раз в sendingInterval
	c.AddJob(time.Duration(*reportIntervalSeconds)*time.Second, func() {
		//отправляем на сервер
		err := sendStorageMetrics(store, "http://"+*addr)
		if err != nil {
			fmt.Println("Попытка отправить метрики завершилась с  ошибками:")
			fmt.Print(err)
		}
		// Сбрасываем PollCount
		// #TODO: в таком случае нужно проверить, что счетчик реально отправился,
		//		а не просто, или нам пофигу?)
		dropCounter(store, metricPollCount)
	})

	c.Serve(200 * time.Millisecond)

}

func configure() {
	flag.Parse()
	loadEnv()
}
