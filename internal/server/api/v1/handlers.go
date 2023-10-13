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

func UnwrapURLParams(handler func(ctx context.Context, params urlParams) (out string, err error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		api.SetContentType(w, api.TypeTextPlain)

		params := getURLParams(r)
		out, err := handler(r.Context(), params) // todo мы почти тут пришли к какому-то универсальному обработчику, черт, типа напрмер типа Metric
		if err != nil {                          // todo вот этот код встречается в соседних обертках
			if err == api.ErrorNotFoundMetric {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "Не найдена метрика %v с именем %v", params.mtype, params.id)
				return
			}
			if err == api.ErrorUnknownMetricType {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Неизвестный тип счетчика: %v", params.mtype)
				return
			} else if err == api.ErrorParseError {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу пропарсить новое значение счетчика типа %v с именем %v значением %v", params.mtype, params.id, params.value)
				return
			}
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Неизвестная ошибка обработки счетчика типа %v с именем %v значением %v: %v", params.mtype, params.id, params.value, err)
			return

		}
		w.Write([]byte(out))
	}
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
	mtype string //type
	id    string
	value string
}

// getURLParams достает из URL маршрута параметры счетчика, такие как
// тип, имя, значение, и возвращает в виде структуры urlParams
func getURLParams(r *http.Request) urlParams {
	return urlParams{
		mtype: chi.URLParam(r, "type"),
		id:    chi.URLParam(r, "name"),
		value: chi.URLParam(r, "value"),
	}
}

func (i API) updateCounterString(ctx context.Context, name string, s string) (int64, error) {
	delta, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, api.ErrorParseError
	}
	return i.store.IncrementCounter(ctx, name, delta)
}

func (i API) updateGaugeString(ctx context.Context, name string, s string) (float64, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, api.ErrorParseError
	}
	return i.store.UpdateGauge(ctx, name, v)
}

// UpdateString обновляет значение в хранилище, имея значение в формате строки. Сам
// просматриваем тип счетчика и решает куда писать
func (i API) UpdateString(ctx context.Context, params urlParams) (out string, err error) {
	switch params.mtype {
	case "counter":
		_, err := i.updateCounterString(ctx, params.id, params.value)
		return "", err
	case "gauge":
		_, err := i.updateGaugeString(ctx, params.id, params.value)
		return "", err
	default:
		return "", api.ErrorUnknownMetricType
	}

}

// GetString обновляет значение в хранилище, имея значение в формате строки. Сам
// просматриваем тип счетчика и решает куда писать.
//
// параметр s не используется, и нужен только для соответствия интерфейсу.
func (i API) GetString(ctx context.Context, params urlParams) (value string, err error) {

	switch params.mtype {
	case "counter":
		c, err := i.store.Counter(ctx, params.id)
		return strconv.FormatInt(c, 10), err // лишний вызов форматирования конечно, но это для редких случаев ошики
	case "gauge":
		g, err := i.store.Gauge(ctx, params.id)
		return strconv.FormatFloat(g, 'f', -1, 64), err
	default:
		return "", api.ErrorUnknownMetricType
		// мне конечно очень не хочется проверять все эти статусы, но с другой стороны это редкие случаи все, то есть замедление
		// на интроспецию ошибок будет идти в редких случаях, когда не тот тип метрики или неправильное значение передано. В идеале это
		// большая редкость
	}

}
