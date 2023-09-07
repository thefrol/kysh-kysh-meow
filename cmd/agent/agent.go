package main

import (
	"fmt"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/scheduler"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

const server = "http://localhost:8080"
const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

var store storage.Storager

func init() {
	store = storage.New()
}

func main() {
	// запуск планировщика
	c := scheduler.New()
	//собираем данные раз в pollingInterval
	c.AddJob(pollInterval, func() {
		fetchMemStats(store)
		fetchAdditionalStats(store)
		// Увеличиваем PollCount
		incrementCounter(store, metricPollCount)
	})
	// отправляем данные раз в sendingInterval
	c.AddJob(reportInterval, func() {
		err := sendStorageMetrics(store, server)
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
