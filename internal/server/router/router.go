package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	apiv1 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v1"
	apiv2 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v2"
	"github.com/thefrol/kysh-kysh-meow/internal/server/middleware"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// apiV1 создает маршруты в роутере, отвечает
// за наследственный API, это когда передача и получение
// значений шли через длинные адреса, например:
// /update/counter/TestCounter/10 или /value/gauge/Alloc
//
// Испрользуется так
// router := chi.NewRouter()
// router.Group(apiV1)
func InstallAPIV1(r chi.Router, v1 apiv1.API) {
	r.Group(func(r chi.Router) {
		// в какой-то момент, когда починят тесты, тут можно будет снять комменты
		//r.With(chimiddleware.AllowContentType("text/plain"))
		r.Get("/value/{type}/{name}", v1.GetValue)

		r.Post("/update/{type:counter}/{name}/{value}", v1.UpdateCounter)
		r.Post("/update/{type:gauge}/{name}/{value}", v1.UpdateGauge)
		r.Post("/update/{type}/{name}/{value}", apiv1.ErrorUnknownType)

	})
}

// apiV2 создает маршруты апи нового образца,
// /update и /value, принимающие джейсон запросы
//
// Испрользуется так
// router := chi.NewRouter()
// router.Group(apiV2)
func apiV2(r chi.Router) {
	r.With(chimiddleware.AllowContentType("application/json"))

	r.Post("/value", apiv2.ValueWithJSON)
	r.Post("/value/", apiv2.ValueWithJSON)
	r.Post("/update", apiv2.UpdateWithJSON)
	r.Post("/update/", apiv2.UpdateWithJSON)
	// как не дублировать маршруты я пока варианта не нашел:
	// если в конце поставить слеш, то без слеша не работает
	// а вроде даже в тестах и так и так иногда бывает
}

// MeowRouter - основной роутер сервера, он отвечает за все мидлвари
// и все маршруты, и даже чтобы на ответы типа 404 и 400 отправлять
// стилизованные ответы.
//
// на входе получает store - объект хранилища, из которого он будет
// брать все нужные данные о метриках
func MeowRouter(store storage.Storager) (router chi.Router) {

	router = chi.NewRouter()

	// настраиваем мидлвари, логгер, распаковщик и запаковщик
	router.Use(middleware.MeowLogging())
	router.Use(middleware.UnGZIP)
	router.Use(middleware.GZIP(middleware.GZIPDefault))

	// Добавляем маршруты, которые я разделил на два раздела
	v1 := apiv1.New(store)
	apiv2.SetStore(store)

	InstallAPIV1(router, v1)
	router.Group(apiV2)

	// а ещё вот HTML страничка, которая тоже по сути относится к apiV1
	// она не объединяется с остальными, потому что не требует
	// application/json или text/plain в заголовках
	router.Get("/", v1.ListMetrics)

	// Тут добавляем стилизованные под кошки-мышки ответы сервера при 404 и 400
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
