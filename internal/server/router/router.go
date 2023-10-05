package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	apiv1 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v1"
	apiv2 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v2"
	"github.com/thefrol/kysh-kysh-meow/internal/server/middleware"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// todo hlog.FromRequest(r).Info() !!!

func MeowRouter(store storage.Storager) (router chi.Router) {

	apiv1.SetStore(store)
	apiv2.SetStore(store)

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

	router.Get("/", apiv1.ListMetrics)
	router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", apiv1.MetricAsURL(apiv1.GetValue))
		r.Post("/", apiv2.ValueWithJSON)
	}) // todo как-то поработать с allowContentType
	router.Route("/update", func(r chi.Router) {
		//r.Use(chimiddleware.AllowContentType("text/plain"))

		r.
			//With(chimiddleware.AllowContentType("application/json")).
			Post("/", apiv2.UpdateWithJSON)
		r.
			Post("/{type:counter}/{name}/{value}", apiv1.MetricAsURL(apiv1.UpdateCounter))
		r.
			Post("/{type:gauge}/{name}/{value}", apiv1.MetricAsURL(apiv1.UpdateGauge))
		r.
			Post("/{type}/{name}/{value}", apiv1.MetricAsURL(apiv1.UpdateUnknownType)) // todo ERROR видно, что хендлер вызывает ошибку
	})

	// как еще вариант могут быть классы, то есть у нас было бы
	// handlers.V1().SetStorage(s)
	// r.Post("/", handlers.V1().Update)
	// - или -
	// r.Post("/", handlers.ByUri.Update)

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
