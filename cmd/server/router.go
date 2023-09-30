package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/cmd/server/middleware"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// todo hlog.FromRequest(r).Info() !!!

func MeowRouter() (router chi.Router) {
	router = chi.NewRouter()

	router.Use(middleware.MeowLogging())
	router.Use(middleware.UnGZIP) // еще бы сообщать, что такой энкодинг не поддерживается
	router.Use(middleware.GZIP(
		middleware.GZIPBestCompression,
		middleware.ContentTypes("text/plain", "text/html", "application/json", "application/xml"),
		middleware.StatusCodes(http.StatusOK),
		middleware.MinLenght(1),
		// todo
		//
		// Хочу чтобы это выглядело так compress.ContentType(...)
	))
	//todo
	//
	// по красоте было бы, если миддлеварь выглядела так router.Use(middleware.Uncompress(Gzip,Deflate,Brotli))
	// и может uncompress.Brotli, uncompress.GZIP
	//
	// а compress вот так compress.GZIP(WithMinLenght(100),WithContentType("text/html")

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

	// for admin
	router.Get("/save/{filename}", func(w http.ResponseWriter, r *http.Request) {
		fname := chi.URLParam(r, "filename")
		err := store.(storage.MemStore).ToFile(fname) // todo
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ololog.Info().Msg("saved")
		w.Write([]byte("saved"))
	})
	router.Get("/load/{filename}", func(w http.ResponseWriter, r *http.Request) {
		fname := chi.URLParam(r, "filename")
		var err error
		store, err = storage.FromFile(fname) // todo
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ololog.Info().Msg("loaded")
		w.Write([]byte("loaded"))
	})

	return router
}
