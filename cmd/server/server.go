// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	//Создать хранилище
	store = storage.New()
}

func main() {
	configure()

	fmt.Printf("^.^ Мяу, сервер работает по адресу %v!\n", *addr)
	err := http.ListenAndServe(*addr, MeowRouter())
	if err != nil {
		fmt.Printf("^0^ не могу запустить сервер: %v \n", err)
	}
}

func configure() {
	flag.Parse()
	loadEnv()
}
