package apiv3

import (
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

type API struct {
	app *app.App
}

func New(app *app.App) API {
	return API{app: app}
}

func (i API) CheckConnection(w http.ResponseWriter, r *http.Request) {
	err := i.app.CheckConnection(r.Context())
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует")
	}
}
