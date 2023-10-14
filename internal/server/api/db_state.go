package api

import (
	"net/http"
)

func CheckConnection(store Operator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := store.Check(r.Context())
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "^0^ Соединение с базой данных отсуствует: %v", err)
		}
	}

}
