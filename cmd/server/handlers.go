package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	//Создать хранилище
	store = storage.New()
}

// updateCounter отвечает за маршрут, по которому будет обновляться счетчик типа counter
// иначе говоря за URL вида: /update/counter/<name>/<value>
// приходящее значение: int64
// поведение: складывать с предыдущим значением, если оно известно
func updateCounter(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseInt(params.Value(), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.Value(), value)
	}
	old, _ := store.Counter(params.Name())
	// по сути нам не надо обрабатывать случай, если значение небыло установлено. Оно ноль, прибавим новое значение и все четко
	new := old + storage.Counter(value)
	store.SetCounter(params.Name(), new)

	fmt.Fprintf(w, "^.^ мур! Меняем Counter %v на %v. Новое значение %v", params.Name(), value, new)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться метрика типа gauge
// иначе говоря за URL вида: /update/gauge/<name>/<value>
// приходящее значение: float64
// поведение: устанавливать новое значение
func updateGauge(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseFloat(params.Value(), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.Value(), value)
	}
	store.SetGauge(params.Name(), storage.Gauge(value))

	fmt.Fprintf(w, "^.^ мур! меняем Gauge %v на %v.", params.Name(), value)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться счетчик неизвестного типа
// без разбора возвращаем 400(Bad Request)
func updateUnknownType(w http.ResponseWriter, r *http.Request, params URLParams) {
	w.WriteHeader(http.StatusBadRequest) //обязательно за вызова w.Write() #INSIGHT добавить в README
	io.WriteString(w, "Фшшш! Я не знаю такой тип счетчика")
	w.Header().Add("Content-Type", "text/plain")
	fmt.Printf("400(BadRequest) at request to %v\n", r.URL.Path)
}

// updateMetricFunc это типа функций обработчков, таких как updateCounter, updateGauge
type updateMetricFunc func(http.ResponseWriter, *http.Request, URLParams)

// makeHandler оборачивает обработчик(например updateCounter) в HandlerFunc
// Проверяет, чтобы марштрут выглядел как надо и заодно парсит его и передает
// в функцию обработчик updateHandleFunc
func makeHandler(fn updateMetricFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// проверки можно отправить в makeHandler
		if r.Method != http.MethodPost {
			//можно использовать http.NotFound
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Мяу! Мы поддерживаем только POST-запросы")
			fmt.Printf("GET request at %v\n", r.URL.Path)
			return
		}
		// Этот блок закомментирован, чтобы пройти автотесты. В тестах отправляется запрос без text/plain и это вызывает ошибку
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Printf("Wront content type at %v\n", r.URL.Path)
		// 	io.WriteString(w, "Мяу! Мы поддерживаем только Content-Type:text/plain")
		// 	return
		// }

		urlparams, err := ParseURL(r.URL.Path)
		if err != nil {
			fmt.Printf("Cant match url %v\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fn(w, r, urlparams)
	}
}
