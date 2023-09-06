// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("^.^ Мяу, это сервер!")
	http.ListenAndServe(":8080", MeowRouter())
}
