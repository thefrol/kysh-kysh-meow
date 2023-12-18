package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (a *ForQuery) GetCounter(w http.ResponseWriter, r *http.Request) {
	var (
		id = chi.URLParam(r, "id")
	)

	v, err := a.Registry.Counter(r.Context(), id)
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
