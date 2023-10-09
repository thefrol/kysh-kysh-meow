package api

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// httpErrorWithLogging отправляет сообщение об ошибке, параллельно дублируя в журнал. Работает быстрее, чем просто две функции отдельно.
// Во-первых, конкатенация строк происходит при помощи Spfrintf, а не сложением, а во вторых один раз на два вызова: и логгера, и http.Error()
//
// w - responseWriter вашего HTTP хендлера
// statusCode - код ответа сервера, напр. 200, 400, http.StatusNotFound, http.StatusOK
// format, params - типичные параметры, как в функции Printf
func HTTPErrorWithLogging(w http.ResponseWriter, statusCode int, format string, params ...interface{}) {
	s := fmt.Sprintf(format, params...)
	log.Error().Str("location", "json update handler").Msg(s)
	http.Error(w, s, statusCode)
	// TODO
	//
	// Возможно это пока единственный повод держать кастомный логгер, чтобы в нем была функция типа withHttpError(w)
}