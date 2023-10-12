package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	apiv1 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v1"
	apiv2 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v2"
	apiv3 "github.com/thefrol/kysh-kysh-meow/internal/server/api/v3"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/middleware"
)

// apiV1 создает маршруты в роутере, отвечает
// за наследственный API, это когда передача и получение
// значений шли через длинные адреса, например:
// /update/counter/TestCounter/10 или /value/gauge/Alloc
//
// Испрользуется так
// router := chi.NewRouter()
// InstallAPIV1(router, store)
func InstallAPIV1(r chi.Router, v1 apiv1.API) {
	r.Group(func(r chi.Router) {
		// в какой-то момент, когда починят тесты, тут можно будет снять комменты
		//r.With(chimiddleware.AllowContentType("text/plain"))
		r.Get("/value/{type}/{name}", apiv1.TextWrapper(v1.GetString))

		r.Post("/update/{type}/{name}/{value}", apiv1.TextWrapper(v1.UpdateString))

	})
}

// InstallAPIV2 создает маршруты апи нового образца,
// /update и /value, принимающие джейсон запросы
//
// Испрользуется так
// router := chi.NewRouter()
// InstallAPIV2(router, store)
func InstallAPIV2(r chi.Router, v2 apiv2.API) {
	r.Group(func(r chi.Router) {
		r.With(chimiddleware.AllowContentType("application/json"))

		update := apiv2.MarshallUnmarshallMerica(v2.UpdateStorage)
		get := apiv2.MarshallUnmarshallMerica(v2.GetStorage)

		r.Post("/value", get)
		r.Post("/value/", get)
		r.Post("/update", update)
		r.Post("/update/", update)

		// как не дублировать маршруты я пока варианта не нашел:
		// если в конце поставить слеш, то без слеша не работает
		// а вроде даже в тестах и так и так иногда бывает
	})

}

// MeowRouter - основной роутер сервера, он отвечает за все мидлвари
// и все маршруты, и даже чтобы на ответы типа 404 и 400 отправлять
// стилизованные ответы.
//
// на входе получает store - объект хранилища, из которого он будет
// брать все нужные данные о метриках
func MeowRouter(store api.Storager, app *app.App) (router chi.Router) {

	router = chi.NewRouter()

	// настраиваем мидлвари, логгер, распаковщик и запаковщик
	router.Use(middleware.MeowLogging())
	router.Use(middleware.UnGZIP)
	router.Use(middleware.GZIP(middleware.GZIPDefault))

	// Добавляем маршруты, которые я разделил на два раздела
	v1 := apiv1.New(store)
	v2 := apiv2.New(store)
	v3 := apiv3.New(app)

	InstallAPIV1(router, v1)
	InstallAPIV2(router, v2)

	//создаем маршрут для проверки соединения с БД
	router.Get("/ping", v3.CheckConnection)

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
