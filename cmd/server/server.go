// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики

package main

import (
	"fmt"
	"io"
	"net/http"
)

type updateHandleFunc func(http.ResponseWriter, *http.Request, URLParams)

// makeHandler оборачивает функцию обработчик для маршрута
// Проверяет, чтобы марштрут выглядел как надо и заодно парсит его и передает
// в функцию обработчик updateHandleFunc
func makeHandler(fn updateHandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlparams, err := ParseUrl(r.URL.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, *urlparams)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request, params URLParams) {
	// проверки можно отправить в makeHandler
	if r.Method != http.MethodPost {
		//можно использовать http.NotFound
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
	mux.Handle("/update/", makeHandler(updateHandler))
}

func main() {

	fmt.Println("Привет, это сервер!")
	http.ListenAndServe(":8080", mux)
}
