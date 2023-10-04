package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/handlers"
	"github.com/thefrol/kysh-kysh-meow/internal/server/middleware"
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

	router.Get("/", handlers.ListMetrics)
	router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.MetricAsURL(handlers.GetValue))
		r.Post("/", handlers.ValueWithJSON)
	}) // todo как-то поработать с allowContentType
	router.Route("/update", func(r chi.Router) {
		//r.Use(chimiddleware.AllowContentType("text/plain"))

		r.
			//With(chimiddleware.AllowContentType("application/json")).
			Post("/", handlers.UpdateWithJSON)
		r.
			Post("/{type:counter}/{name}/{value}", handlers.MetricAsURL(handlers.UpdateCounter))
		r.
			Post("/{type:gauge}/{name}/{value}", handlers.MetricAsURL(handlers.UpdateGauge))
		r.
			Post("/{type}/{name}/{value}", handlers.MetricAsURL(handlers.UpdateUnknownType)) // todo ERROR видно, что хендлер вызывает ошибку
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
