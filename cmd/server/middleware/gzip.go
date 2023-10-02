package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

// GZIP это мидлварь для сервера, которая сжимает содержимое запроса
// и оформляет все заголовки. Имеет множество настроек, которые передаются
// при создании, такие как:
//
//	MinLenght(длинна) - минимальная длинна тела в байтах, чтобы сжать тело запроса
//	ContentType(заголовок1, заголовок2, ... ) разрешенные к передаче заголовки, обычно text/plain,text/html, application/json
//	StatusCodes(код1, код2, ...) разрешенные к сжиманию по кодам ответов, обычно 200
//	BestCompession, BestSpeed - уровень компрессии
//
// В общем виде, на сервере chi создание мидлвари выглядит следующим образом,
// router.Use(GZIP(BestCompresion,MinLenght(50),ContentTypes("text/plain"),StatusCodes(http.StatusOK)))
func GZIP(funcOpts ...gzipFuncOpt) func(http.Handler) http.Handler {

	// Тут описываемся сама мидлварь, получившая opts
	// в качестве настроек
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Eсли клиент поддерживает gzip, то подменяем врайтер, не забыв его закрыть, и отправляем
			// запрос дальше по цепочке
			if acceptsEncoding(r, "gzip") {
				cw := NewCompressedWriter(w, funcOpts...)
				w = cw
				defer cw.Close()
			}
			next.ServeHTTP(w, r)

		})
	}
}

// CompressedWriter это обертка над ResponseWriter,
// которая если нужно, будет использовать сжимать данные
// а если не нужно, просто отправит как есть
//
// Стоит обязательно создавать конструктором NewCompressedWriter,
// иначе будут пробелмы с статус кодом по умолчанию
//
// Вся концепция этого класса держится на том, в http.ResponseWrite
// устроен следующим образом: к тому моменту, как впервые используется
// Write, все заголовки и код ответа уже должны быть записаны, и
// впоследствии не могут быть изменены. А это мы можем должаться этого первого
// Write() и решить, будем ли мы сжимать данные или нет.
//
// Например, если мы отвечаем 404 или 400, то сжимать не надо
// И надо сжимать, если Content-Type: text или application/json
//
// Главная проблема, конечно с минимальной длинной, после которой начинается сжатие,
// это означает что первые байты мы должны записывать в какой-то буфер и если он переполнится, то
// активировать сжатие, и все писать потом уже в сжимальщик gzip.Writer
type CompressedWriter struct {
	opts gzipOptions

	archiveWriter  io.WriteCloser
	originalWriter http.ResponseWriter
	statusCode     int

	fistBytesBuffer bytes.Buffer // сюда будут записаны первые minLenght байт
	bytesWritten    int

	ignore  bool // true - не будем сжимать
	checked bool // true- прошел все проверки по другим параментрам кроме минамальной длинны и собирается быть сжат

}

func NewCompressedWriter(originalWriter http.ResponseWriter, funcOpts ...gzipFuncOpt) *CompressedWriter {

	// Создаем структуру настроек для GZIP,
	// после применения к ней опцефункций,
	// она будет передана в замыкание и будет уже работать с мидлварью
	opts := gzipOptions{
		CompressionLevel: gzip.BestCompression,
	}
	for _, f := range funcOpts {
		f(&opts)
	}
	return &CompressedWriter{
		opts:           opts,
		originalWriter: originalWriter,
		statusCode:     200, // сразу ставим StatusOK, потому что такой код у нас по умолчанию
	}

	//todo я думаю нам нужны какие-то настройки по умолчанию, заголовки и статус коды нарзерешенные
}

func (cw CompressedWriter) Header() http.Header {
	return cw.originalWriter.Header()
}

func (cw *CompressedWriter) Write(bb []byte) (int, error) {

	// если уже было выбрано, что компрессию игнорируем, то просто пишем в обернурый врайтер
	if cw.ignore {
		return cw.originalWriter.Write(bb)
	}

	// если архиватор уже создан, то тупо пишем в него и не думаем
	if cw.archiveWriter != nil {
		return cw.archiveWriter.Write(bb)
	}
	//todo мы можем из Content-Lenght прочитать сколько у меня байт будет

	// тут мы проверяем, что этот ответ будет сжат, если проходит проверки, то устанавливаем checked=true
	// стоит отменить, что checked и ignore не взаимозаменяемы. Во время первой записи у нас ни тот ни другой
	// не установлены, и checked в основном используется, чтобы сэкономить пару тактов на проверку, чтобы
	// не проверять кучу хедеров во время каждой записи, хотя может это и лишняя забота
	if !cw.checked {

		// до этой отметки никакой записи в тело originalWriter, тот что мы обернули,
		// не происходило, а ниже уже начинаются первые записи, значит надо подготовиться
		// и записать все что потом уже не запишется, например код ответа

		// либо мы прошли проверки на архивацию checked=true, либо архивания отменилась ignore==true
		if cw.opts.CheckStatusCode(cw.statusCode) && cw.opts.CheckContentTypes(cw.originalWriter.Header()) {
			cw.checked = true

		} else {
			// мы не будем использовать компрессию, устанавливаем флаг ignore и записываем байты в обернутый врайтер
			cw.ignore = true
			cw.confirmStatusCode()
			return cw.originalWriter.Write(bb)
		}
	}

	// нам прошел буфер размером len, если с учетом буфера нам все же не хватает байт
	// чтобы активировать архивацию, то записываем в буфер, и возвращаем управление
	//
	// TODOs
	//
	// из исходного кода chi.middleware мы знаем, что можно просто прочитать content-length
	// и сразу узнаем, сколько у нас там байт,
	// а коли контент-ленгтх не указан, то напрямую компрессируем
	if cw.bytesWritten+len(bb) < cw.opts.MinLenght {
		n, err := cw.fistBytesBuffer.Write(bb)
		cw.bytesWritten += n
		return n, err
	}

	// Теперь мы прошли минимальный порог байтов, значит архивируем:
	// создадим архиватор и сбросить все что уже накопилось

	gz, err := gzip.NewWriterLevel(cw.originalWriter, cw.opts.CompressionLevel)
	if err != nil {
		return 0, fmt.Errorf("не могу создать gzip.Writer: %v", err)
	}

	// не забудем подсказать в заголовках, что мы архивируем содержание
	cw.originalWriter.Header().Add("Content-Encoding", "gzip")
	cw.originalWriter.Header().Del("Content-Length")
	cw.confirmStatusCode()

	cw.archiveWriter = gz

	// если во временный буфер уже успели что-то записать то сбрасываем в архиватор
	if cw.bytesWritten > 0 {
		n, err := io.Copy(cw.archiveWriter, &cw.fistBytesBuffer)
		if err != nil {
			return int(n), err
		}

	}

	// todo
	//
	// надо бы проверять, что он уже не сжат,

	nn, err := gz.Write(bb)
	if err != nil {
		ololog.Error().Str("Location", "compression middleware").Msg("error writing to first bytes buffer when writing, should not happen")
		return nn, err
		// todo тут ещё можно скипнуть с архива и попробовать например без архивации, если что-то идет не так. Например может можно вызвать recover?
	}
	cw.bytesWritten += nn

	return nn, nil

	// todo
	//
	// 1. пул архиваторов отдельным пакетом. Или пул-пакет) пул-пулак
	// пул-пулак даже написать хочется
	// короче создаются динамически, поддерживается число десять,
	// может быть и больше, но тогда они будут утилизироваться,
	// или больше нельзя)
	//
	// 2. архивация идет по порядку, надо это учитывать, что может быть заархивировано в каком-то порядке
}

func (cw *CompressedWriter) WriteHeader(status int) {
	// Запоминаем код ответа и используем его
	// перед первым Write() в оригинальный
	// responseWriter
	cw.statusCode = status
}

// confirmStatusCode подтвержает полученный статус код, и записывает его
// в оригинальный врайтер. В отличие от WriteHeader(), который просто запоминает код
// эта функция  имеет внутреннее значение для мидлвари,  она
// уже реально записывает статус код на выход.
//
// Тут мы проверяем, а не пришел ли нам статус код 0
// по сути это тот же код 200, просто связанный с ошибками на других уровнях обработки. На это стоит обратить
// пристальное внимание, но такая ошибка не должна валить сервер
func (cw *CompressedWriter) confirmStatusCode() {
	if cw.statusCode == 0 {
		ololog.Error().Str("location", "server/middleware/gzip").Msg("Кто-то записал мне код статуса 0, вместо 200, это где-то после меня случилось. Я поменял статус код на 200, но нужно обязательно проверить, что-то работает не так")
		cw.statusCode = 200
	}
	if cw.statusCode != 200 {
		cw.originalWriter.WriteHeader(cw.statusCode)
	}
}

func (cw *CompressedWriter) Close() {
	// Если порог минимального количество байт так и не пройдет, то просто сбрасываем все в
	// обернутый врайтер
	if cw.archiveWriter == nil {
		cw.confirmStatusCode()
		io.Copy(cw.originalWriter, &cw.fistBytesBuffer)
	} else {
		// а если архиватор создан, то его надо закрыть
		cw.archiveWriter.Close()
	}
	// это уже не часть интерфейса, но нужна для нашей мидвари, чтобы закрыть gzip.Writer
}

var _ http.ResponseWriter = (*CompressedWriter)(nil)

// acceptsEncoding возвращает true, если на основании запроса r
// можно сказать, клиент поддерживает кодирование в encoding
//
// if acceptsEncoding(r, "gzip"){...}
func acceptsEncoding(r *http.Request, encoding string) bool {
	for _, v := range r.Header.Values("Accept-Encoding") {
		if strings.Contains(v, encoding) {
			return true
		}
	}
	return false
}
