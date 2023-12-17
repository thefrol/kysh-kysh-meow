package queryhandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (a *API) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	var (
		id    = chi.URLParam(r, "id")
		delta = chi.URLParam(r, "delta")
	)

	d, err := strconv.ParseInt(delta, 10, 64)
	if err != nil {
		//todo нужен модуль httperror )
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "handler: IncrementCounter() не могу пропарсить инкремент для счетчика: %v", err)
		return
	}

	v, err := a.Registry.IncrementCounter(r.Context(), id, d)
	if err != nil {
		if errors.Is(err, domain.ErrorMetricNotFound) {
			api.HTTPErrorWithLogging(w, http.StatusNotFound, "handler: GetCounter() не найдена метрика : %v", err)
			return
		}
		api.HTTPErrorWithLogging(w, http.StatusBadRequest, "handler: GetCounter() handler : %v", err)
		return
	}

	api.SetContentType(w, api.TypeTextPlain)
	w.Write([]byte(strconv.FormatInt(v, 10)))
}
