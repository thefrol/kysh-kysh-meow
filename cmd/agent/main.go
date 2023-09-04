package main

import (
	"fmt"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

const server = "http://localhost:8080"
const (
	pollingInterval = 2 * time.Second
	sendingInterval = 10 * time.Second
)

var store storage.Storager

func init() {
	store = storage.New()
}

func main() {
	//проверка метрик
	lost, exclude, err := parseMetrics()
	if err != nil {
		fmt.Printf("Проверка метрик выдала ошибку: %v", err)
		return
	}
	if len(lost) > 0 {
		fmt.Printf("Невозможно получить следующие необходимые метрики памяти: %v", lost)
		return
	}
	fmt.Printf("Необязательные метрики памяти будут проигнорированы: %v\n", exclude)

	// запуск планировщика
	c := Chron{}
	//собираем данные раз в pollingInterval
	c.AddJob(pollingInterval, func() {
		saveMemStats(store, exclude)
		saveAdditionalStats(store)
		updateCounter(store)
	})
	// отправляем данные раз в sendingInterval
	c.AddJob(sendingInterval, func() {
		err := sendStorageMetrics(store, server)
		if err != nil {
			fmt.Println("Попытка отправить метрики завершилось с  ошибками:")
			fmt.Print(err)
		}
	})

	c.Serve(200 * time.Millisecond)

}
