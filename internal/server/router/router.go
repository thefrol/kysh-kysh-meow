package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	handler "github.com/thefrol/kysh-kysh-meow/internal/server/handlers"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"

	"github.com/thefrol/kysh-kysh-meow/internal/server/middleware"
)

const (
	CompressionTreshold  = 50
	CompressionBufferLen = 2048
)

// MeowRouter - основной роутер сервера, он отвечает за все мидлвари
// и все маршруты, и даже чтобы на ответы типа 404 и 400 отправлять
// стилизованные ответы.
//
// на входе получает store - объект хранилища, операции над которым он будет проворачивать
func MeowRouter(store api.Operator, key string) (router chi.Router) {

	router = chi.NewRouter()

	// создадим уровень приложения
	m := manager.Registry{
		Counters: &storage.CounterAdapter{Op: store},
		Gauges:   &storage.GaugeAdapter{Op: store},
	}

	// настраиваем мидлвари, логгер, распаковщик и запаковщик
	router.Use(middleware.MeowLogging())
	if key != "" {
		router.Use(
			middleware.CheckSignature(key),
			middleware.SignResponse(key))
	}
	router.Use(middleware.UnGZIP)
	router.Use(middleware.GZIP(CompressionTreshold, CompressionBufferLen))

	// Первая часть маршрутов - из первых спринтов. Тут используются
	// параметры из маршрута для установки значений метрик, за эти маршруты
	// отвечает структура query
	router.Group(func(r chi.Router) {
		// хендлеры сгрупированы в эту структурку, тут все что надо
		// для этих самых простых хендлеров
		query := handler.ForQuery{
			Registry: m,
		}

		// на данный момент тесты неправильно работают с контент-тайпом,
		// например они не присылают правильный контент-тайп в некоторых случаях,
		// а если так-то пока мы не можем ставить эту проверку
		// но когда тесты починят, это стоит сделать //todo
		//
		// r.With(chimiddleware.AllowContentType("text/plain"))

		// Настроим хендлеры получения метрик. Мы добавляем ещё один
		// маршрут для метрик с невалидным типом, при обращении к этому
		// маршруту, мы всегда возвращаем 400 Bad Request
		r.Get("/value/counter/{id}", query.GetCounter)
		r.Get("/value/gauge/{id}", query.GetGauge)
		r.Get("/value/{unknown}/{id}", BadRequest)

		// Настроим хендлеры изменения метрик. Аналогино
		// поступаем с невалидным маршрутом, который
		// всегда вернет 400 Bad Request
		r.Post("/update/counter/{id}/{delta}", query.IncrementCounter)
		r.Post("/update/gauge/{id}/{value}", query.UpdateGauge)
		r.Post("/update/{unknown}/{id}/{value}", BadRequest)
	})

	// Создаем маршруты для обработки JSON запросов
	router.Group(func(r chi.Router) {
		//r.With(chimiddleware.AllowContentType("application/json"))

		r.Post("/value", api.HandleJSONRequest(api.Retry3Times(store.Get)))
		r.Post("/value/", api.HandleJSONRequest(api.Retry3Times(store.Get)))
		r.Post("/update", api.HandleJSONRequest(api.Retry3Times(store.Update)))
		r.Post("/update/", api.HandleJSONRequest(api.Retry3Times(store.Update)))
		r.Post("/updates", api.HandleJSONBatch(api.Retry3Times(store.Update)))
		r.Post("/updates/", api.HandleJSONBatch(api.Retry3Times(store.Update)))

		// TODO
		//
		// подозрительно похоже на абстракную фабрику
		// update := api.HandleJSONRequest(store.Update)
		// get := api.HandleJSONRequest(store.Get)
		//
		//
		// как не дублировать маршруты я пока варианта не нашел:
		// если в конце поставить слеш, то без слеша не работает
		// а вроде даже в тестах и так и так иногда бывает
	})

	// Создаем маршруты, показывающие статус сервера. Страница со всемми метриками,
	// и пинг БД
	router.Group(func(r chi.Router) {
		router.Get("/ping", api.PingStore(store))
		router.Get("/", api.DisplayHTML(store))
	})

	// Тут добавляем стилизованные под кошки-мышки ответы сервера при 404 и 400,
	// Кроме того, мы подменяем MethodNotAllowed на NotFound
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

// BadRequest это специальный хендлер, который возвращает ошибку 400 Bad Request
func BadRequest(w http.ResponseWriter, r *http.Request) {
	api.HTTPErrorWithLogging(w,
		http.StatusBadRequest,
		"0-0 ошибка в запросе")
}
