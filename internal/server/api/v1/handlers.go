// Этот пакет содержит хендлеры старого образца, где мы передавали значения при помощи URL
// например, /update/counter/PollCount/120
package apiv1

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

// API это колленция http.HanlderFunc, которые обращаются к единому хранилищу store
type API struct {
	store api.Storager
}

// New создает новую
func New(store api.Storager) API {
	if store == nil {
		panic("Хранилище - пустой указатель")
	}
	return API{store: store}
}

// updateCounter отвечает за маршрут, по которому будет обновляться счетчик типа counter
// иначе говоря за URL вида: /update/counter/<name>/<value>
// приходящее значение: int64
// поведение: складывать с предыдущим значением, если оно известно
func (i API) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	api.SetContentType(w, api.TypeTextPlain)

	value, err := strconv.ParseInt(params.value, 10, 64)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу пропарсить новое значение counter: %v", err)
		return
	}

	_, err = i.store.UpdateCounter(r.Context(), params.name, value)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "не могу обновить счетчик: %v", err)
		return
	}
}

// updateGauge отвечает за маршрут, по которому будет обновляться метрика типа gauge
// иначе говоря за URL вида: /update/gauge/<name>/<value>
// приходящее значение: float64
// поведение: устанавливать новое значение
func (i API) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	api.SetContentType(w, api.TypeTextPlain)

	value, err := strconv.ParseFloat(params.value, 64)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу пропарсить новое значение gauge: %v", err)
		return
	}
	_, err = i.store.UpdateGauge(r.Context(), params.name, value)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "не могу обновить gauge: %v", err)
		return
	}
}

// ErrorUnknownType возвращает клиенту ошибку 400(Bad Request)
// и отправляет информационное сообщение, о том, что сервер не знает такого типа счетчика
func ErrorUnknownType(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	api.SetContentType(w, api.TypeTextPlain)

	api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "неизвестный тип счетчика: %v", params.metric)
}

// getValue возвращает значение уже записанной метрики,
// если метрика ранее не была записана, возвращает http.StatusNotFound
// если попытка обратиться к метрике несуществующего типа http.StatusNotFound
func (i API) GetValue(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	var (
		value string
		found bool
		err   error
	)

	// Здесь мы получаем значение метрики с полученным именем,
	//
	// TODO
	//
	// Меня не отпускает ощущение, что можно чу-у-у-у-уточку ускорить
	// весь этот алгоритм, если проверять не строку целиком, а
	// например первую буковку
	// - если "c" - перед нами counter
	// - если "g" - gauge
	//
	// все проверки мы уже сделали заранее!
	switch params.metric {
	case metrica.CounterName:
		var v int64
		v, found, err = i.store.Counter(r.Context(), params.name)
		value = fmt.Sprint(v)
	case metrica.GaugeName:
		var v float64
		v, found, err = i.store.Gauge(r.Context(), params.name)
		value = fmt.Sprint(v)
		// TODO
		//
		// Этот кусок кода, конечно тоже малопонятен, я вообще за такие ункции в хранилишще
		// SetCounterString(), CounterString(), возможно даже SetString(type, name, sval), GetString(type, name, sval)
	}

	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка получения метрики: %v", err)
		return
	}
	if !found {
		api.HTTPErrorWithLogging(w, http.StatusNotFound, "метрика [%v]%v не найдена в хранилище", params.metric, params.name)
		return
	}

	w.Write([]byte(value))
}

// listMetrics выводит список всех известных на данный момент метрик
func (i API) ListMetrics(w http.ResponseWriter, r *http.Request) {
	api.SetContentType(w, api.TypeTextHTML)

	htmlTemplate := `
	{{ if .ListCounters}}
		Counters
			<ul> 
				{{range .ListCounters -}}
					<li>{{.}}</li>
				{{ end }}
			</ul>
	{{ else }}
		<p>No Counters</p>
	{{end}}
	{{ if .ListGauges}}
		Gauges
			<ul> 
				{{range .ListGauges -}}
					<li>{{.}}</li>
				{{ end }}
			</ul>
	{{ else }}
		<p>No Gauges</p>
	{{end}}
	`

	tmpl, err := template.New("simple").Parse(htmlTemplate)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не удалось пропарсить HTML шаблон: %v", err)
		return
	}
	err = tmpl.Execute(w, i.store)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка запуска шаблона HTML: %v", err)
		return
	}
}

type urlParams struct {
	metric string //type
	name   string
	value  string
}

// getURLParams достает из URL маршрута параметры счетчика, такие как
// тип, имя, значение, и возвращает в виде структуры urlParams
func getURLParams(r *http.Request) urlParams {
	return urlParams{
		metric: chi.URLParam(r, "type"),
		name:   chi.URLParam(r, "name"),
		value:  chi.URLParam(r, "value"),
	}
}
