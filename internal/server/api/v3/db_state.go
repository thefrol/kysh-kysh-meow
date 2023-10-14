package apiv3

import (
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

func CheckConnection(store api.Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := store.Check(r.Context())
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует: %v", err)
		}
	}

}
