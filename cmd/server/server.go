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

func main() {
	fmt.Println("^.^ Мяу, это сервер!")
	http.ListenAndServe(":8080", MeowRouter())
}
