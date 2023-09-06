// Сервер Мяу-мяу
// Умеет сохранять и передавать такие метрики: counter, gauge

package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var router chi.Router

func init() {
	router = chi.NewRouter()
	router.Get("/value/{type}/{name}", makeHandler(getMetric))
	router.Route("/update", func(r chi.Router) {
		r.Post("/{type:counter}/{name}/{value}", makeHandler(updateCounter))
		r.Post("/{type:gauge}/{name}/{value}", makeHandler(updateGauge))
		r.Post("/{type}/{name}/{value}", makeHandler(updateUnknownType))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("^0^ оуууоо! такой метод или андрес недоступен"))
	})

}

func main() {
	fmt.Println("^.^ Мяу, это сервер!")
	http.ListenAndServe(":8080", router)
}
