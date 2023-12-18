package queryhandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

// GetGauge это хендлер, который возвращает метрику типа gauge,
// с идентификатором id
func (a *API) GetGauge(w http.ResponseWriter, r *http.Request) {

	// достанем из query идентификатор имя нашей метрики
	var (
		id = chi.URLParam(r, "id")
	)

	// обратимся к слою приложения и получим из хранилища
	// значение нашей метрики
	v, err := a.Registry.Gauge(r.Context(), id)
	if err != nil {
		// если метрика не найдена, то мы пишем в ответ статус 404
		if errors.Is(err, domain.ErrorMetricNotFound) {
			api.HTTPErrorWithLogging(w,
				http.StatusNotFound,
				"handler: GetGauge() не найдена метрика %v: %v", id, err)
			return
		}

		// в остальных случаях подразумеваем 400 - плохой запрос
		api.HTTPErrorWithLogging(w,
			http.StatusBadRequest,
			"handler: GetGauge() ошибка для метрики %v: %v", id, err)
		return
	}

	// если метрика получена, то ставим контент тайп
	// и пишем прям в тело ответа значение
	api.SetContentType(w, api.TypeTextPlain)
	w.Write([]byte(strconv.FormatFloat(v, 'g', -1, 64))) // todo эта функция могла бы быть частью domain или моделей, чего-то такого более общего
}
