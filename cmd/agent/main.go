package main

import (
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

const server = "http://localhost:8080"

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
	//основной цикл
	saveMemStats(store, exclude)
	err = sendStorageMetrics(store, server)
	if err != nil {
		fmt.Println("Попытка отправить метрики завершилось с  ошибками:")
		fmt.Print(err)
	}
	fmt.Println(store.Gauge("Alloc"))
	fmt.Println(store.Gauge("HeapSys"))
}
