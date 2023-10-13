package apiv3

import (
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func CheckConnection(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.CheckConnection(r.Context())
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует")
		}
	}

}
