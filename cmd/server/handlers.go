package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

// updateWithJSON обновляет значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica. Для счетчиков типа counter исползует поле delta и прибавляет к
// текущему значению, для счетчиков типа gauge заменяет текущее значение новым из поля Value.
// В ответ записывает структуру metrica.Metrica с обновленным значением
func updateWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		ololog.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "^0^ не могу размаршалить тело сообщения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = m.Validate()
	if err != nil {
		ololog.Error().Str("location", "json update handler(validate metrica)").Msgf("Получена невалидная структура: %+v", err)
		http.Error(w, "Полученная структура оформлена неправильно: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case metrica.CounterName:
		oldVal, _ := store.Counter(m.ID)
		val := oldVal + metrica.Counter(*m.Delta)
		store.SetCounter(m.ID, val) // todo: we need update counter!!!
	case metrica.GaugeName:
		store.SetGauge(m.ID, metrica.Gauge(*m.Value)) // todo точно надо в структере записать как gauge и counter поля
	default:
		ololog.Error().Str("location", "json update handler").Msgf("Cant update metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	// returning можно сделать кой-то отдельной функицей
	var result metrica.Metrica
	var found bool
	switch m.MType {
	case metrica.CounterName:
		var c metrica.Counter
		c, found = store.Counter(m.ID) //todo меня бесит это возвращающее в два параметра, надо сделать функцию Lookup
		result = c.Metrica(m.ID)
	case metrica.GaugeName:
		var g metrica.Gauge
		g, found = store.Gauge(m.ID)
		result = g.Metrica(m.ID)
	default:
		ololog.Error().Str("location", "json update handler(on return)").Msgf("Cant update metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	if !found {
		// todo такое ощущение, что нужен какой-то слой логики, везде одно и то же делаю как будто
		http.Error(w, "returned counter not found", http.StatusBadRequest)
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&result, w)
	if err != nil {
		ololog.Error().Str("location", "json update handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

// value WithJSON возвращает значение счетчика JSON запросом. Читает из запроса тело в
// формате metrica.Metrica, у которого должны быть заполенены поля MType и ID,
// TДля счетчиков типа counter записывает значнием в поле delta,
// для счетчиков типа gauge в поле value. В ответ
// записывает структуру metrica.Metrica с обновленным значением
func valueWithJSON(w http.ResponseWriter, r *http.Request) {
	//todo mix of two handlers + save+ response
	m := metrica.Metrica{}
	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		ololog.Error().Str("location", "json update handler").Msg("Cant unmarshal data in body")
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var result metrica.Metrica
	var found bool
	switch m.MType {
	case metrica.CounterName:
		var c metrica.Counter
		c, found = store.Counter(m.ID) //todo меня бесит это возвращающее в два параметра, надо сделать функцию Lookup
		result = c.Metrica(m.ID)
	case metrica.GaugeName:
		var g metrica.Gauge
		g, found = store.Gauge(m.ID)
		result = g.Metrica(m.ID)
	default:
		ololog.Error().Str("location", "json value handler(on return)").Msgf("Cant get valye metric with type %v, no such metric", m.MType)
		http.Error(w, "unsupported type", http.StatusBadRequest)
		return
	}

	if !found {
		// todo такое ощущение, что нужен какой-то слой логики, везде одно и то же делаю как будто
		http.Error(w, "metrica not found", http.StatusNotFound)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(&result, w)
	if err != nil {
		ololog.Error().Str("location", "json get handler(on return)").Msgf("Cant marshal a return for %v %v", m.MType, m.ID)
		http.Error(w, "cant marshal result", http.StatusInternalServerError)
		return
	}
}

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
	case metrica.CounterName:
		value, found = store.Counter(params.name)
	case metrica.GaugeName:
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
	w.Header().Add("Content-Type", "text/html")

	// todo: сделать это в html-разметке при помощи template
	const indent = "    "

	cl := store.ListCounters()
	gl := store.ListGauges()

	if len(cl)+len(gl) == 0 {
		w.Write([]byte("No metrics stored"))
		return
	}
	if len(cl) > 0 {
		fmt.Fprintln(w, "Counters:")
		for _, v := range cl {
			fmt.Fprintln(w, indent, v)
		}

	}
	if len(gl) > 0 {
		fmt.Fprintln(w, "Gauges:")
		for _, v := range gl {
			fmt.Fprintln(w, indent, v)
		}
	}
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
