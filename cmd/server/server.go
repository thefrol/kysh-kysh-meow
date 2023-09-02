// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"fmt"
	"net/http"
)

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.Handle("/update/counter/", makeHandler(updateCounter))
	mux.Handle("/update/gauge/", makeHandler(updateGauge))
	mux.Handle("/update/", makeHandler(updateUnknownType))
}

func main() {
	fmt.Println("^.^ Мяу, это сервер!")
	http.ListenAndServe(":8080", mux)
}
