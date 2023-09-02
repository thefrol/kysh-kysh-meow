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
		// проверки можно отправить в makeHandler
		if r.Method != http.MethodPost {
			//можно использовать http.NotFound
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Мяу! Мы поддерживаем только POST-запросы")
			fmt.Printf("GET request at %v\n", r.URL.Path)
			return
		}
		// Пройти автотесты
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Printf("Wront content type at %v\n", r.URL.Path)
		// 	io.WriteString(w, "Мяу! Мы поддерживаем только Content-Type:text/plain")
		// 	return
		// }

		urlparams, err := ParseUrl(r.URL.Path)
		if err != nil {
			fmt.Printf("Cant match url %v\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fn(w, r, *urlparams)
	}
}

var mux *http.ServeMux

func init() {
	//так можно и вообще роутинг убрать
	mux = http.NewServeMux()
	mux.Handle("/update/counter/", makeHandler(updateCounter))
	mux.Handle("/update/gauge/", makeHandler(updateGauge))
	mux.Handle("/update/", makeHandler(updateUnknownType))
}

func main() {

	fmt.Println("Привет, это сервер!")
	http.ListenAndServe(":8080", mux)
}
