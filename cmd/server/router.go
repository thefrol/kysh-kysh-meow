package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// todo hlog.FromRequest(r).Info() !!!

func MeowRouter() (router chi.Router) {
	router = chi.NewRouter()

	router.Use(MeowLogging())

	router.Get("/", listMetrics)
	router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", makeHandler(getValue))
		r.Post("/", valueWithJSON)
	}) // todo как-то поработать с allowContentType
	router.Route("/update", func(r chi.Router) {
		//r.Use(chimiddleware.AllowContentType("text/plain"))

		r.
			//With(chimiddleware.AllowContentType("application/json")).
			Post("/", updateWithJSON)
		r.
			Post("/{type:counter}/{name}/{value}", makeHandler(updateCounter))
		r.
			Post("/{type:gauge}/{name}/{value}", makeHandler(updateGauge))
		r.
			Post("/{type}/{name}/{value}", makeHandler(updateUnknownType))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("^0^ оуууоо! такой метод не дотупен по этому адресу"))
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(404)
		w.Write([]byte("^0^ оуууоо! Нет такой страницы"))
	})

	return router
}
