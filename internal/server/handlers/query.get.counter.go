package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
	"github.com/thefrol/kysh-kysh-meow/internal/server/httpio"
)

func (a *ForQuery) GetCounter(w http.ResponseWriter, r *http.Request) {
	var (
		id = chi.URLParam(r, "id")
	)

	v, err := a.Registry.Counter(r.Context(), id)
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
