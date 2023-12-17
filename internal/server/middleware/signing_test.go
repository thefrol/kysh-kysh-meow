package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
)

// Handler простая ручка для тестового сервера запросов
// перед ней будет стоять мидлварь, которая при ошибке
// подписи тупо не пропустит до сюда
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Response"))
}

// тест проверяет работу подписывающей мидлвари.
// Она должна подписать и прочитать запрос
func TestSigning(t *testing.T) {
	// Ключ шифрования
	var (
		key = "123"
	)

	// Запрос
	var (
		body        = "data_for_signing"
		requestSign = "CyHw9JM8fq5pvtsVbGVPYq/+qXahhUCyl18O85/DFt8="
	)

	// Ответ
	var (
		statusCode   = http.StatusOK
		responseSign = "XdNsy/vA+t7X+zIELkzdoHbvXylCa3F2vxBlf0y3kQo="
	)

	//создадим наш сервер
	route := "/path"

	s := chi.NewRouter()
	s.Use(CheckSignature(key))
	s.Get(route, Handler)

	// создадим запрос
	r := httptest.NewRequest(http.MethodGet, route, bytes.NewBuffer([]byte(body)))
	r.Header.Set(sign.SignHeaderName, requestSign)

	// запускаем запрос
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	// обрабатываем ответ
	res := w.Result()

	assert.Equal(t, statusCode, res.StatusCode)
	assert.Equal(t, responseSign, res.Header.Get(sign.SignHeaderName))
}

// тест проверяет работу подписывающей мидлвари.
// Приходящий запрос имеем неправильную подпись
func Test_Signing_BadSign(t *testing.T) {
	// Ключ шифрования
	var (
		key = "123"
	)

	// Запрос
	//
	// тут подпись запроса не валидная
	var (
		body        = "data_for_signing"
		requestSign = "bad_sign_CyHw9JM8fq5pvtsVbGVPYq/+qXahhUCyl18O85/DFt8="
	)

	// Ответ
	var (
		statusCode = http.StatusBadRequest
	)

	//создадим наш сервер
	route := "/path"

	s := chi.NewRouter()
	s.Use(CheckSignature(key), SignResponse(key))
	s.Get(route, Handler)

	// создадим запрос
	r := httptest.NewRequest(http.MethodGet, route, bytes.NewBuffer([]byte(body)))
	r.Header.Set(sign.SignHeaderName, requestSign)

	// запускаем запрос
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	// обрабатываем ответ
	res := w.Result()

	assert.Equal(t, statusCode, res.StatusCode)
}
