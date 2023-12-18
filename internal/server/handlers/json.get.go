package handlers

import (
	"errors"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (h *ForJSON) Get(w http.ResponseWriter, r *http.Request) {
	var m metrica.Metrica

	err := easyjson.UnmarshalFromReader(r.Body, &m)
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не могу размаршалить в json тело запроса: %v", err)
		return
	}

	// todo validate
	// and domain.MEtrica или прям тут

	switch m.MType {
	case "counter":
		v, err := h.Registry.Counter(r.Context(), m.ID)
		if err != nil {
			if errors.Is(err, domain.ErrorMetricNotFound) {
				api.HTTPErrorWithLogging(w,
					http.StatusNotFound,
					"Не удалось найти счетчик %v: %v", m.ID, err)
				return
			}
			api.HTTPErrorWithLogging(w,
				http.StatusInternalServerError,
				"Не могу получить значение cчетчика %v: %v", m.ID, err)
			return
		}
		m.Delta = &v
	case "gauge":
		v, err := h.Registry.Gauge(r.Context(), m.ID)
		if err != nil {
			if errors.Is(err, domain.ErrorMetricNotFound) {
				api.HTTPErrorWithLogging(w,
					http.StatusNotFound,
					"Не удалось найти гауж %v: %v", m.ID, err)
				return
			}
			api.HTTPErrorWithLogging(w,
				http.StatusInternalServerError,
				"Не могу получить значение гаужа %v: %v", m.ID, err)
			return
		}
		m.Value = &v
	default:
		api.HTTPErrorWithLogging(w,
			http.StatusBadRequest,
			"Неизвестный тип метрики %v, запрос с именем %v", m.MType, m.ID)
		return
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(m, w)
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не замаршалить ответ: %v", err)
		return
	}
}
