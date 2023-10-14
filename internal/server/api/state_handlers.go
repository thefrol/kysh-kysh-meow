package api

import (
	"net/http"
	"text/template"
)

// PingStore создает хендлер, который пингует хранилище. Если хранилище может связаться с базой данных,
// то оно ответит 200(OK), иначе 500(Internal Server Error). Отвечает без ошибки, только если хранилищем установлена
// база данных, и если с ней есть связь.
func PingStore(store Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := store.Check(r.Context())
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует: %v", err)
		}
	}

}

// DisplayHTML создает хендлер HTTP запроса. Этот хендлер формирует простую
// HTML страничку, где указаны все известные на данный момент метрики
func DisplayHTML(op Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SetContentType(w, TypeTextHTML)

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
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не удалось пропарсить HTML шаблон: %v", err)
			return
		}
		cs, gs, err := op.List(r.Context())
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка получения списка метрик из хранилища: %v", err)
			return
		}
		err = tmpl.Execute(w, struct {
			ListCounters []string
			ListGauges   []string
		}{ListCounters: cs, ListGauges: gs})
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "Ошибка запуска шаблона HTML: %v", err)
			return
		}
	}
}
