package handlers

import (
	"html/template"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
)

const htmlTemplate = `
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

type ForHTML struct {
	Labels scan.Labels
}

func (html ForHTML) Dashboard(w http.ResponseWriter, r *http.Request) {
	api.SetContentType(w, api.TypeTextHTML)

	tmpl, err := template.New("simple").Parse(htmlTemplate)
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не удалось пропарсить HTML шаблон: %v", err)
		return
	}

	labels, err := html.Labels.Get(r.Context())
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Ошибка получения списка метрик из хранилища: %v", err)
		return
	}

	data := struct {
		ListCounters []string
		ListGauges   []string
	}{
		ListCounters: labels["counters"],
		ListGauges:   labels["gauges"],
	}

	err = tmpl.Execute(w, data) // todo Это все как-то может просто мапу обрабатывать
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Ошибка запуска шаблона HTML: %v", err)
		return
	}
}
