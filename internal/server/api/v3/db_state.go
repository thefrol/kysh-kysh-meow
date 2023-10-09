package apiv3

import (
	"context"
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

const waitTime = time.Second * 1

type API struct {
	app *app.App
}

func New(app *app.App) API {
	return API{app: app}
}

func (i API) CheckConnection(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(i.app.Context(), waitTime)
	defer cancel()
	err := i.app.CheckConnection(ctx)
	if err != nil {
		api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует")
	}
}
