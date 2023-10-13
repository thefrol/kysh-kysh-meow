// Этот пакет содержит хендлеры старого образца, где мы передавали значения при помощи URL
// например, /update/counter/PollCount/120
package apiv1

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

func UnwrapURLParams(handler api.Operation) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		api.SetContentType(w, api.TypeTextPlain)

		in := getURLParams(r).Parse()

		// Валидиуем полученную структуру
		if in.ID == "" {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена направильно заполенная струкура %+v: имя метрики не может быть пустым", in)
			return
		}

		arr, err := handler(r.Context(), in)

		if err != nil { // todo вот этот код встречается в соседних обертках
			if err == api.ErrorDeltaEmpty || err == api.ErrorValueEmpty {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка входных данных: %v", err)
			} else if err == api.ErrorNotFoundMetric {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "Не найдена метрика %v с именем %v", in.MType, in.ID)
				return
			} else if err == api.ErrorUnknownMetricType {
				api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Неизвестный тип счетчика: %v", in.MType)
				return
			}
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Неизвестная ошибка обработки счетчика типа %v с именем %v значением %v: %v", in.MType, in.ID, in.Value, err)
			return

		}

		// поскольку мы обрабаываем кучей, то как бы нужно взять из массива одно какое-то
		// возможно мне понадобится еще один дополнительный оберточник
		if len(arr) != 1 {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "После обработки операции над хранилищем получено неправильное количество выходящих значений")
			return
		}
		out := arr[0]

		s, err := ValueString(out)
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка конвертации значения метрики %v в строку: %v", out.MType, err)
		}

		w.Write([]byte(s))

		// TODO
		//
		// В общем мы видим как выглядит любая такая оболочка
		// запрос -> конвертер в metrica -> opetation -> конвертирование обратно -> ответ
		//
		// В целом, даже проверка ошибок +- одинаковая + ещё устаовка контент тайпа, да
		//
		// Возможно конвертер это такой класс типа) класс класс на классе и классом погоняет
	}
}

// listMetrics выводит список всех известных на данный момент метрик
func ListMetrics(op api.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		cs, gs, err := op.List(r.Context())
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
}

type urlParams struct {
	mtype string //type
	id    string
	value string
}

func (p urlParams) Parse() metrica.Metrica {
	out := metrica.Metrica{}
	out.MType = p.mtype
	out.ID = p.id

	// если в значении не пустая строка, то пытаемся распарсить,
	// просто и инт и флоат, чтобы не проверять типы счетчиков лишний раз

	if p.value != "" {
		c, err := strconv.ParseInt(p.value, 10, 64)
		if err == nil {
			out.Delta = &c
		}
		g, err := strconv.ParseFloat(p.value, 64)
		if err == nil {
			out.Value = &g
		}
	}

	return out

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

func ValueString(m metrica.Metrica) (string, error) {

	switch m.MType {
	case "counter":
		if m.Delta == nil {
			return "", api.ErrorDeltaEmpty
		}
		return strconv.FormatInt(*m.Delta, 10), nil // лишний вызов форматирования конечно, но это для редких случаев ошики
	case "gauge":
		if m.Value == nil {
			return "", api.ErrorValueEmpty
		}
		return strconv.FormatFloat(*m.Value, 'f', -1, 64), nil
	default:
		return "", api.ErrorUnknownMetricType
		// мне конечно очень не хочется проверять все эти статусы, но с другой стороны это редкие случаи все, то есть замедление
		// на интроспецию ошибок будет идти в редких случаях, когда не тот тип метрики или неправильное значение передано. В идеале это
		// большая редкость
	}

}
