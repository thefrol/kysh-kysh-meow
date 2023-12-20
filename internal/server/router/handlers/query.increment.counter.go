package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
)

func (a *ForQuery) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	var (
		id    = chi.URLParam(r, "id")
		delta = chi.URLParam(r, "delta")
	)

	d, err := strconv.ParseInt(delta, 10, 64)
	if err != nil {
		//todo нужен модуль httperror )
		httpio.HTTPErrorWithLogging(w, http.StatusBadRequest, "handler: IncrementCounter() не могу пропарсить инкремент для счетчика: %v", err)
		return
	}

	v, err := a.Registry.IncrementCounter(r.Context(), id, d)
	if err != nil {
		if errors.Is(err, domain.ErrorMetricNotFound) {
			httpio.HTTPErrorWithLogging(w, http.StatusNotFound, "handler: GetCounter() не найдена метрика : %v", err)
			return
		}
		httpio.HTTPErrorWithLogging(w, http.StatusBadRequest, "handler: GetCounter() handler : %v", err)
		return
	}

	w.Header().Set("Content-Type", httpio.TypeTextPlain)
	w.Write([]byte(strconv.FormatInt(v, 10)))
}
