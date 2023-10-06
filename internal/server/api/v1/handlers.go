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
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func SetStore(s storage.Storager) {
	store = s
}

// updateCounter отвечает за маршрут, по которому будет обновляться счетчик типа counter
// иначе говоря за URL вида: /update/counter/<name>/<value>
// приходящее значение: int64
// поведение: складывать с предыдущим значением, если оно известно
func UpdateCounter(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

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
func UpdateGauge(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	value, err := strconv.ParseFloat(params.value, 64)
	if err != nil {
		w.Header().Add("Content-Type", "text/plain")
		http.Error(w, "^0^ Ошибка значения, не могу пропарсить", http.StatusBadRequest)
		return
	}
	store.SetGauge(params.name, metrica.Gauge(value))
	w.Header().Add("Content-Type", "text/plain")
}

// ErrorUnknownType возвращает клиенту ошибку 400(Bad Request)
// и отправляет информационное сообщение, о том, что сервер не знает такого типа счетчика
func ErrorUnknownType(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	http.Error(w, "Фшшш! Я не знаю такой тип счетчика", http.StatusBadRequest)
}

// getValue возвращает значение уже записанной метрики,
// если метрика ранее не была записана, возвращает http.StatusNotFound
// если попытка обратиться к метрике несуществующего типа http.StatusNotFound
func GetValue(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

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
func ListMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

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
		ololog.Error().Str("location", "server/handlers/listCounters").Err(err).Msg("Не удается создать и пропарсить HTML шаблон")
		http.Error(w, "error creating template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, store)
	if err != nil {
		ololog.Error().Str("location", "server/handlers/listCounters").Err(err).Msg("Ошибка запуска шаблона HTML")
		http.Error(w, "error executing template", http.StatusInternalServerError)
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