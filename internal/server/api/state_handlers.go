package api

import (
	"net/http"
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
