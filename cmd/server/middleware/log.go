package middleware

import (
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

func MeowLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			d := time.Duration(0)
			wr := wrapWriter(w)
			d = countTime(func() {
				// запустить обработку
				next.ServeHTTP(wr, r)
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
			ololog.Info().
				Str("method", r.Method).
				Str("uri", r.RequestURI).
				Dur("duration", d).
				Msg("Request ->")

			ololog.Info().
				Int("statusCode", wr.statusCode).
				Int("size", wr.bytesWritten).
				Msg("Response ->")
		})
	}
}

// wrapWriter Оборачивает оригинальный http.ResponseWriter оболочкой writeWrapper, которая хранит
// некоторые данные о манипуляциях с этим врайтером. Например, какой статус код был записан,
// сколько данных записано в байтах
func wrapWriter(w http.ResponseWriter) *writerWrapper {
	return &writerWrapper{originalWriter: w, statusCode: 200}
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

// writeWrapper это оболочка для интерфейса http.ResponseWriter, которая хранит некоторые параметры о
// манипуляциях с ней: сколько было записино, и какие данные под конец, так же отслеживает чтобы итоговый статус код был правильно записан
type writerWrapper struct {
	originalWriter http.ResponseWriter
	bytesWritten   int
	statusCode     int
}

func (ww *writerWrapper) Header() http.Header {
	return ww.originalWriter.Header()
}
func (ww *writerWrapper) Write(bb []byte) (int, error) {
	n, err := ww.originalWriter.Write(bb)
	ww.bytesWritten += n
	return n, err
}

func (ww *writerWrapper) WriteHeader(statusCode int) {
	if ww.bytesWritten > 0 {
		// я кое-что узнал в перерыве, что после использования Write()
		// заголовки нельзя больше переписать даже при помощи WriteHeader()
		// поэтому проверяем

		// TODO: проверки я бы вынес в отдельную мидлварь тогда,
		// например, ещё можно проверять какой контент тайп мы выдаем
		ololog.Error().Str("location", "http processing").Msg("Попытка записи заголовков после использования функции Write(). Заголовки и статус уже не изменить")
	}
	ww.originalWriter.WriteHeader(statusCode)
	ww.statusCode = statusCode
}

// проверить что writeWrapper отвечает интерфейсу http.ResponseWriter
var _ http.ResponseWriter = (*writerWrapper)(nil)

// тут по-хорошему сделать бы инъекцию зависимости и предоставить нужный журнал, но нормального интерфейса у Олологгера у нас нет
