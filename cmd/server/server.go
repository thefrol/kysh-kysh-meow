// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"fmt"
	"net/http"

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

	fmt.Printf("^.^ Мяу, сервер работает по адресу %v!\n", cfg.Addr)
	err := http.ListenAndServe(cfg.Addr, MeowRouter())
	if err != nil {
		fmt.Printf("^0^ не могу запустить сервер: %v \n", err)
	}
}
