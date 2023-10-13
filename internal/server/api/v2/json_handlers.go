// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package apiv2

import (
	"context"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

// API это колленция http.HanlderFunc, которые обращаются к единому хранилищу store
type API struct {
	store api.Operator
}

// New создает новую апиху для джейсонов
func New(store api.Operator) API {
	if store == nil {
		panic("Хранилище - пустой указатель")
	}
	return API{store: store}
}

// MarshallUnmarshallMetrica - это функция обертка, которая размаршаливает и замаршаливает значения полученные по HTTP в структуру metrica.Metrica,
// используется, чтобы избавиться от дублирования кода в конктретных хендлерах /value и /update
func MarshallUnmarshallMerica(handler func(context.Context, ...metrica.Metrica) (out []metrica.Metrica, err error)) func(http.ResponseWriter, *http.Request) {

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
		arr, err := handler(r.Context(), in)
		if err != nil {
			if err == api.ErrorNotFoundMetric {
				api.HTTPErrorWithLogging(w, http.StatusNotFound, "В хранилище не найдена метрика %v", in.ID)
				return
			}
			api.HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка обновления метрики %+v: %v", in, err) // todo, а как бы сделать так, чтобы %v подсвечивался
			return
		}

		// поскольку мы обрабаываем кучей, то как бы нужно взять из массива одно какое-то
		// возможно мне понадобится еще один дополнительный оберточник
		if len(arr) != 1 {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "После обработки операции над хранилищем получено неправильное количество выходящих значений")
			return
		}
		out := arr[0]

		// Замаршаливаем результат работы хендлера
		_, _, err = easyjson.MarshalToHTTPResponseWriter(&out, w)
		if err != nil {
			api.HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не могу замаршалить выходной джейсон: %v", err)
			return
		}
	}

}
