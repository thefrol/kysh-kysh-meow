// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package apiv2

import (
	"context"
	"errors"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

var (
	ErrorDeltaEmpty = errors.New("Поле Delta не может быть пустым, для когда id=counter")
	ErrorValueEmpty = errors.New("Поле Value не может быть пустым, для когда id=gauge")
)

// API это колленция http.HanlderFunc, которые обращаются к единому хранилищу store
type API struct {
	store api.Storager
}

// New создает новую
func New(store api.Storager) API {
	if store == nil {
		panic("Хранилище - пустой указатель")
	}
	return API{store: store}
}

func MarshallUnmarshallMerica(handler func(context.Context, metrica.Metrica) (out metrica.Metrica, err error)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		/*
			1. Размаршаливаем полученное сообщение в структуру metrica.Metrica
			2. Валидиуем полученную структуру
			3. Запускаем обработчик handler(), он как-то там работает с хранилищем, или ещё чего, что нам не очень важно.
				И возвращает структуру out.
			4. Замаршаливаем результат работы хендлера
		*/

		// Размаршаливаем полученное сообщение в структуру metrica.Metrica
		in := metrica.Metrica{}
		err := easyjson.UnmarshalFromReader(r.Body, &in)
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу размаршалить тело сообщения: %v", err)
			return
		}
		defer r.Body.Close()

		// Валидиуем полученную структуру
		if in.ID == "" {
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена направильно заполенная струкура %+v: имя метрики не может быть пустым", in)
			return
		}

		// Запускаем обработчик handler()
		out, err := handler(r.Context(), in)
		if err != nil {
			if err == api.ErrorNotFoundMetric {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "В хранилище не найдена метрика %v", in.ID)
				return
			}
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка обновления метрики %+v: %v", in, err) // todo, а как бы сделать так, чтобы %v подсвечивался
			return
		}

		// Замаршаливаем результат работы хендлера
		_, _, err = easyjson.MarshalToHTTPResponseWriter(&out, w)
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не могу замаршалить выходной джейсон: %v", err)
			return
		}
	}

}

func (a API) UpdateStorage(ctx context.Context, in metrica.Metrica) (newVal metrica.Metrica, err error) {

	switch in.MType {
	case "counter":
		if in.Delta == nil {
			return empty, ErrorDeltaEmpty
		}
		c, err := a.store.IncrementCounter(ctx, in.ID, *in.Delta)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Delta: &c}, err // это получается отправится в хип
	case "gauge":
		if in.Value == nil {
			return empty, ErrorValueEmpty
		}
		g, err := a.store.UpdateGauge(ctx, in.ID, *in.Value)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Value: &g}, err
	default:
		return empty, api.ErrorUnknownMetricType
	}
}

func (a API) GetStorage(ctx context.Context, in metrica.Metrica) (newVal metrica.Metrica, err error) {
	switch in.MType {
	case "counter":
		c, err := a.store.Counter(ctx, in.ID)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Delta: &c}, err // это получается отправится в хип
	case "gauge":
		g, err := a.store.Gauge(ctx, in.ID)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Value: &g}, err
	default:
		return empty, api.ErrorUnknownMetricType
	}
}

var empty = metrica.Metrica{}
