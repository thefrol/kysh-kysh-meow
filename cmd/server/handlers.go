package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
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
		http.Error(w, "^0^ Ошибка значения, не могу пропарсить", http.StatusBadRequest)
		return
	}
	old, _ := store.Counter(params.name)
	// по сути нам не надо обрабатывать случай, если значение небыло установлено. Оно ноль, прибавим новое значение и все четко
	new := old + metrica.Counter(value)
	store.SetCounter(params.name, new)
	w.Header().Add("Content-Type", "text/plain")
	ololog.Info().Str("location", "handlers").Fields(params).Msgf("^.^ мур! Меняем Counter %v на %v. Новое значение %v", params.name, value, new)
}

// updateGauge отвечает за маршрут, по которому будет обновляться метрика типа gauge
// иначе говоря за URL вида: /update/gauge/<name>/<value>
// приходящее значение: float64
// поведение: устанавливать новое значение
func updateGauge(w http.ResponseWriter, r *http.Request, params URLParams) {
	value, err := strconv.ParseFloat(params.value, 64)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, "^0^ Ошибка значения, не могу пропарсить", http.StatusBadRequest)
		return
	}
	store.SetGauge(params.name, metrica.Gauge(value))
	w.Header().Add("Content-Type", "text/plain")
}

// updateGauge отвечает за маршрут, по которому будет обновляться счетчик неизвестного типа
// без разбора возвращаем 400(Bad Request)
// #TODO переименовать в BadRequest
func updateUnknownType(w http.ResponseWriter, r *http.Request, params URLParams) {
	w.Header().Add("Content-Type", "text/plain")
	http.Error(w, "Фшшш! Я не знаю такой тип счетчика", http.StatusBadRequest)
}

// getValue возвращает значение уже записанной метрики,
// если метрика ранее не была записана, возвращает http.StatusNotFound
// если попытка обратиться к метрике несуществующего типа http.StatusNotFound
func getValue(w http.ResponseWriter, r *http.Request, params URLParams) {
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

// listMetrics выводит список всех известных на данный момент метрик
func listMetrics(w http.ResponseWriter, r *http.Request) {
	b := strings.Builder{}
	const indent = "    "

	cl := store.ListCounters()
	gl := store.ListGauges()

	if len(cl)+len(gl) == 0 {
		b.WriteString("No metrics stored")
		w.Write([]byte(b.String()))
		return
	}
	if len(cl) > 0 {
		fmt.Fprintln(&b, "Counters:")
		for _, v := range cl {
			fmt.Fprintln(&b, indent, v)
		}

	}
	if len(gl) > 0 {
		fmt.Fprintln(&b, "Gauges:")
		for _, v := range gl {
			fmt.Fprintln(&b, indent, v)
		}
	}
	w.Write([]byte(b.String()))
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
		fn(w, r, p)
	}
}

type URLParams struct {
	metric string //type
	name   string
	value  string
}
