// Этот пакет содержит хендлеры старого образца, где мы передавали значения при помощи URL
// например, /update/counter/PollCount/120
package apiv1

import (
	"context"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
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
func (i API) UpdatePlain(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	api.SetContentType(w, api.TypeTextPlain)

	err := UpdateString(r.Context(), i.store, params.metric, params.name, params.value)
	if err != nil {

		if err == api.ErrorUnknownMetricType {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Неизвестный тип счетчика: %v", params.metric)
			return
		} else if err == api.ErrorParseError {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу пропарсить новое значение счетчика типа %v с именем %v значением %v", params.metric, params.name, params.value)
			return
		}
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Неизвестная ошибка обновления счетчика типа %v с именем %v значением %v: %v", params.metric, params.name, params.value, err)
		return

	}

}

// getValue возвращает значение уже записанной метрики,
// если метрика ранее не была записана, возвращает http.StatusNotFound
// если попытка обратиться к метрике несуществующего типа http.StatusNotFound
func (i API) GetValue(w http.ResponseWriter, r *http.Request) {
	params := getURLParams(r)

	s, err := GetString(r.Context(), i.store, params.metric, params.name)

	if err != nil {

		if err == api.ErrorNotFoundMetric {
			api.HTTPErrorWithLogging(w, http.StatusNotFound, "Счетчик типа %v с именем %v не найден", params.metric, params.name)
			return
		} else if err == api.ErrorUnknownMetricType {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Неизвестный тип счетчика: %v", params.metric)
			return
		}
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Неизвестная ошибка получения счетчика счетчика типа %v с именем %v значением %v: %v", params.metric, params.name, params.value, err)
		return

	}

	w.Write([]byte(s))
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
	cs, gs, err := i.store.List(r.Context())
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка получения списка метрик из хранилища: %v", err)
		return
	}
	err = tmpl.Execute(w, struct {
		ListCounters []string
		ListGauges   []string
	}{ListCounters: cs, ListGauges: gs})
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

func updateCounterString(ctx context.Context, store api.Storager, name string, s string) (int64, error) {
	delta, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, api.ErrorParseError
	}
	return store.IncrementCounter(ctx, name, delta)
}

func updateGaugeString(ctx context.Context, store api.Storager, name string, s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, api.ErrorParseError
	}
	return store.UpdateGauge(ctx, name, v)
}

// UpdateString обновляет значение в хранилище, имея значение в формате строки. Сам
// просматриваем тип счетчика и решает куда писать
func UpdateString(ctx context.Context, store api.Storager, mtype string, name string, s string) (err error) {
	switch mtype {
	case "counter":
		_, err := updateCounterString(ctx, store, name, s)
		return err
	case "gauge":
		_, err := updateGaugeString(ctx, store, name, s)
		return err
	default:
		return api.ErrorUnknownMetricType
	}

}

// GetString обновляет значение в хранилище, имея значение в формате строки. Сам
// просматриваем тип счетчика и решает куда писать
func GetString(ctx context.Context, store api.Storager, mtype string, name string) (value string, err error) {

	switch mtype {
	case "counter":
		c, err := store.Counter(ctx, name)
		return strconv.FormatInt(c, 10), err // лишний вызов форматирования конечно, но это для редких случаев ошики
	case "gauge":
		g, err := store.Gauge(ctx, name)
		return strconv.FormatFloat(g, 'f', -1, 64), err
	default:
		return "", api.ErrorUnknownMetricType
		// мне конечно очень не хочется проверять все эти статусы, но с другой стороны это редкие случаи все, то есть замедление
		// на интроспецию ошибок будет идти в редких случаях, когда не тот тип метрики или неправильное значение передано. В идеале это
		// большая редкость
	}

}
