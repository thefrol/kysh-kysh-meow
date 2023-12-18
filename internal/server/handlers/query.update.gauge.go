package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (a *ForQuery) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	var (
		id    = chi.URLParam(r, "id")
		value = chi.URLParam(r, "value")
	)

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		// а может это часть пакета log,
		// log.HTTPError(http.StatusOK).Str("type","gauge").Str("id",id).Msg(...)
		// пока не знаю как это реализовать)
		api.HTTPErrorWithLogging(w,
			http.StatusBadRequest,
			"handler: UpdateGauge() не могу пропарсить"+
				" новое значение %v для гаужа %v: %v", value, id, err)
		return
	}

	v, err = a.Registry.UpdateGauge(r.Context(), id, v)
	if err != nil {
		if errors.Is(err, domain.ErrorMetricNotFound) {
			api.HTTPErrorWithLogging(w,
				http.StatusNotFound,
				"handler: UpdateGauge() не найдена гаужа %v: %v", id, err)
			return
		}
		api.HTTPErrorWithLogging(w,
			http.StatusBadRequest,
			"handler: UpdateGauge() не удалось обновить %v на %v : %v", id, value, err)
		return
	}

	api.SetContentType(w, api.TypeTextPlain)
	w.Write([]byte(strconv.FormatFloat(v, 'g', -1, 64)))
}
