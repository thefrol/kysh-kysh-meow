package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// updateCounter отвечает за маршрут, по которому будет обновляться счетчик типа counter
// иначе говоря за URL вида: /update/counter/<name>/<value>
// приходящее значение: int64
// поведение: складывать с предыдущим значением, если оно известно
func updateCounter(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseInt(params.value, 10, 64)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.value, value)
		return
	}
	old, _ := store.Counter(params.name)
	// по сути нам не надо обрабатывать случай, если значение небыло установлено. Оно ноль, прибавим новое значение и все четко
	new := old + metrica.Counter(value)
	store.SetCounter(params.name, new)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "^.^ мур! Меняем Counter %v на %v. Новое значение %v", params.name, value, new)

	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться метрика типа gauge
// иначе говоря за URL вида: /update/gauge/<name>/<value>
// приходящее значение: float64
// поведение: устанавливать новое значение
func updateGauge(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseFloat(params.value, 64)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "^0^ Ошибка значение, не могу пропарсить %v в %T", params.value, value)
		return
	}
	store.SetGauge(params.name, metrica.Gauge(value))
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "^.^ мур! меняем Gauge %v на %v.", params.name, value)

	fmt.Printf("200(OK) at request to %v\n", r.URL.Path)
}

// updateGauge отвечает за маршрут, по которому будет обновляться счетчик неизвестного типа
// без разбора возвращаем 400(Bad Request)
// #TODO переименовать в BadRequest
func updateUnknownType(w http.ResponseWriter, r *http.Request, params URLParams) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest) //обязательно за вызова w.Write() #INSIGHT добавить в README
	io.WriteString(w, "Фшшш! Я не знаю такой тип счетчика")
	fmt.Printf("400(BadRequest) at request to %v\n", r.URL.Path)
}

func getMetric(w http.ResponseWriter, r *http.Request, params URLParams) {
	var value fmt.Stringer
	var found bool
	switch params.metric {
	case "counter":
		value, found = store.Counter(params.name)
	case "gauge":
		value, found = store.Gauge(params.name)
	default:
		http.NotFound(w, r)
		return
	}

	if !found {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(value.String()))
}

// updateMetricFunc это типа функций обработчков, таких как updateCounter, updateGauge
type updateMetricFunc func(http.ResponseWriter, *http.Request, URLParams)

// makeHandler оборачивает обработчик(например updateCounter) в HandlerFunc
// Проверяет, чтобы марштрут выглядел как надо и заодно парсит его и передает
// в функцию обработчик updateHandleFunc
func makeHandler(fn updateMetricFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := URLParams{
			metric: chi.URLParam(r, "type"),
			name:   chi.URLParam(r, "name"),
			value:  chi.URLParam(r, "value"),
		}
		// Этот блок закомментирован, чтобы пройти автотесты. В тестах отправляется запрос без text/plain и это вызывает ошибку
		// if slices.Contains(r.Header("Content-Type"),"text/plain)" { // возможно это все щеё неправильно. Мне приходило такое "text/plain; encoging: utf8"
		// 	w.WriteHeader(http.StatusNotFound)
		// 	fmt.Printf("Wront content type at %v\n", r.URL.Path)
		// 	io.WriteString(w, "Мяу! Мы поддерживаем только Content-Type:text/plain")
		// 	return
		// }
		fn(w, r, p)
	}
}

type URLParams struct {
	metric string //type
	name   string
	value  string
}
