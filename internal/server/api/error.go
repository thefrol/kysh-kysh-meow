package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/lib/retry"
	"github.com/thefrol/kysh-kysh-meow/lib/retry/fails"
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

// Retry3Times позволяет повторить операцию ресколько раз
//
// retriableHandler:=Retry3Times(handler)
// out, err:=retryableHandler(ctx, in)
func Retry3Times(op Operation) Operation {
	return func(ctx context.Context, d ...datastruct) (out []datastruct, err error) {
		err =
			retry.This(
				func() error {
					out, err = op(ctx, d...)
					return err
				},
				retry.If(fails.OnDial),
				retry.Attempts(3),
				retry.DelaySeconds(1, 3, 5, 7),
				retry.OnRetry(
					func(i int, err error) {
						log.Info().Msgf("Выполняется повторная попытка %v после ошибки %v", i, err)
					}),
			)
		if err != nil {
			return nil, err
		}
		return
	}
}
