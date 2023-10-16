// Этот пакет содержит хендлеры нового образца,
// где мы передаем значения при помощи json-запросов
// по маршрутам /update и /value
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// HandleJSONRequest создает HTTP хендлер. Это функция обертка, которая размаршаливает и замаршаливает значения,
// полученные по HTTP в структуру metrica.Metrica, и запускает операцию над хранилищем op
//
// используется, чтобы избавиться от дублирования кода в конктретных хендлерах /value и /update
func HandleJSONBatch(handler func(context.Context, ...metrica.Metrica) (out []metrica.Metrica, err error)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		/*
			1. Размаршаливаем полученное сообщение в структуру metrica.Metrica
			2. Валидиуем полученную структуру
			3. Запускаем обработчик handler(), он как-то там работает с хранилищем, или ещё чего, что нам не очень важно.
				И возвращает структуру out.
			4. Замаршаливаем результат работы хендлера
		*/

		// Размаршаливаем полученное сообщение в структуру metrica.Metrica
		in := make([]metrica.Metrica, 40)
		err := json.NewDecoder(r.Body).Decode(&in) // todo тут использовать easyjson, но надо придумать способ запустить кодогенерацию
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusBadRequest, "Не могу размаршалить тело сообщения: %v", err)
			return
		}
		defer r.Body.Close()

		// Валидиуем полученные структуры
		for _, r := range in {
			if r.ID == "" {
				HTTPErrorWithLogging(w, http.StatusBadRequest, "Получена направильно заполенная струкура %+v: имя метрики не может быть пустым", in)
				return
			}

		}

		// Запускаем обработчик handler()
		out, err := handler(r.Context(), in...)
		if err != nil {
			if errors.Is(err, ErrorNotFoundMetric) {
				HTTPErrorWithLogging(w, http.StatusNotFound, "В хранилище не найдена одна из метрик %v", in)
				return
			}
			HTTPErrorWithLogging(w, http.StatusBadRequest, "Ошибка в работе хендлера метрике %+v: %v", in, err) // todo, а как бы сделать так, чтобы %v подсвечивался
			return
			// todo тут точно надо будет поиграть с обертками
		}

		// Замаршаливаем результат работы хендлера
		err = json.NewEncoder(w).Encode(&out)
		if err != nil {
			HTTPErrorWithLogging(w, http.StatusInternalServerError, "Не могу замаршалить выходной джейсон: %v", err)
			return
		}
	}

}
