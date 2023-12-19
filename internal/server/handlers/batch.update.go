package handlers

import (
	"errors"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (h *ForBatch) Update(w http.ResponseWriter, r *http.Request) {
	var b metrica.Metricas

	err := easyjson.UnmarshalFromReader(r.Body, &b)
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не могу размаршалить в json тело запроса: %v", err)
		return
	}

	// todo
	//
	// в идеале тут нужны какие-то транзации или мютекс хотя бы!
	// Ну или нам так норм. Впринципе, все операции у меня
	// атомарные
	var res metrica.Metricas
	for _, m := range b {
		v, err := h.Manager.UpdateMetrica(r.Context(), m)
		if err != nil {
			if errors.Is(err, domain.ErrorMetricNotFound) {
				api.HTTPErrorWithLogging(w,
					http.StatusNotFound,
					"Не удалось найти метрику %v типа %v: %v", m.ID, m.MType, err)
				return
			}

			api.HTTPErrorWithLogging(w,
				http.StatusBadRequest,
				"Получена неправильная метрика %v.%v: %v", m.MType, m.ID, err)
			return
		}

		res = append(res, v)
	}

	_, _, err = easyjson.MarshalToHTTPResponseWriter(res, w)
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не замаршалить ответ: %v", err)
		return
	}
}
