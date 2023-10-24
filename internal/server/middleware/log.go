package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/lib/intercept"
)

func MeowLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			d := time.Duration(0)
			wr := intercept.New(w, make([]byte, 0, 1024))
			//defer wr.Close() // очень показательная ошибка!
			// тут мы в дефер передали wr со всеми старыми значениями
			// wr.StatusCode и всяких полей.
			// а по логике моей структурки, я записываю этот StatusCode
			// в ответ, приложение потом запишет туда 500
			// но в close() будет звучать все равно исходный 0,
			// который был там записан во время вызова defer
			//
			// тут конечно очень большое поле для выбора как это фиксить
			// можно передавать в дефер замыкание, можно вызывать клоуз
			// по ссылки и тогда передавать по ссылке объект в дефер
			//
			// я выбрал вариант где Close() имеет ресивер по указателю,
			// собсно в этом вся беда и была
			defer wr.Close()

			d = countTime(func() {
				// запустить обработку
				next.ServeHTTP(&wr, r)
			})

			// TODO
			// Возможно в одном ответе мы можем передавать просто две структуры! и в сообщении какие-то самые важные моменты
			//
			// Вообще мне не очень нравится такой формат, как в задании
			// хотелось бы видеть - пришел запрос номер такой-то
			// запрос такой-то
			// тут сообщения от хендера посередине
			// ответ такой-то
			// Чтобы сообщения были как бы обернуты миддлеваром, вот начало вот конец.
			// Или в контексте реквеста передается как-то логгер, куда может писать хендлер и тогда его сообщения
			// будут как-то отдельно форматироваться
			log.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Bool("gzippedRequest", encoded(r, "gzip")).
				Dur("duration", d).
				Msg("Request ->")

			log.Info().
				Int("statusCode", wr.StatusCode()).
				Str("Content-Type", wr.Header().Get("Content-Type")).
				Int("size", wr.Buf().Len()).
				// todo add gzipped response flag
				Msg("Response ->")
		})
	}
}

// countTime засекает время исполнения функции аргумента
func countTime(f func()) (d time.Duration) {
	defer func() {
		d = time.Since(time.Now())
	}()
	f()
	return
	// хрена се, эта функция работает примерно в тысячу раз быстрее чем прошлая с указателем
}
