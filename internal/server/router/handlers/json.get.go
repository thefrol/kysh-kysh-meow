package handlers

import (
	"errors"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
)

func (h *ForJSON) Get(w http.ResponseWriter, r *http.Request) {
	var m metrica.Metrica

	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		httpio.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не могу размаршалить в json тело запроса: %v", err)
		return
	}

	v, err := h.Manager.GetMetrica(r.Context(), m)
	if err != nil {
		if errors.Is(err, app.ErrorMetricNotFound) {
			httpio.HTTPErrorWithLogging(w,
				http.StatusNotFound,
				"Не удалось найти метрику %v типа %v: %v", m.ID, m.MType, err)
			return
		}
		httpio.HTTPErrorWithLogging(w,
			http.StatusBadRequest,
			"Получена неправильная метрика %v.%v: %v", m.MType, m.ID, err)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(v, w)
	if err != nil {
		httpio.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не замаршалить ответ: %v", err)
		return
	}
}
