// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики

package main

import (
	"fmt"
	"io"
	"net/http"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	//переписать бы интерфейс так, чтобы он сразу на входе парсил строчку с командами
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Мяу! Мы поддерживаем только POST-запросы")
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Мяу! Мы поддерживаем только Content-Type:text/plain")
		return
	}
	io.WriteString(w, "^.^ мур!")
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler()) // отвечает на все запросы - не найдено 404
	mux.Handle("/update/", http.HandlerFunc(updateHandler))
}

func main() {

	fmt.Println("Привет, это сервер!")
	http.ListenAndServe(":8080", mux)
}
