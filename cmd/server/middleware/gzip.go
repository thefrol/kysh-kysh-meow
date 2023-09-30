package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/slices"
)

// UnGZIP распаковывает запросы, закодированные при помощи GZIP, и пропускает все
// остальное мимо ушей, отслеживает чтобы хотя бы один из заголовков Content-Encoding
// содержал подстроку gzip
func UnGZIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// пропускаем, поскольку мы не обрабатываем такие сжималки
		if !encoded(r, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// в случае если перед нами gzip, заменяем исходное тело запроса на обертку с gzip
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			ololog.Error().Str("location", "middleware/gzip").Strs("Content-Encoding", r.Header.Values("Content-Encoding")).Err(err).Msg("Cant unzip a request body")
			http.Error(w, "cant unzip", http.StatusBadRequest)
			//todo try recover and send data as is
			return
		}
		defer r.Body.Close()
		r.Body = gz //todo: make a pool of zgips and return only buffer, not gzipeers

		next.ServeHTTP(w, r)

	})
}

type gzipOptions struct {
	MinLenght          int
	contentTypes       []string
	CompressionLevel   int
	allowedStatusCodes []int
}

func (o gzipOptions) CheckContentTypes(h http.Header) bool {
	for _, ct := range o.contentTypes {
		for _, c := range h.Values("Content-Type") {
			if strings.Contains(c, ct) {
				return true
			}
		}
	}
	return false
}

func (o gzipOptions) CheckStatusCode(code int) bool {
	return slices.Contains[int](o.allowedStatusCodes, code)
}

type gzipFuncOpt func(*gzipOptions)

func ContentTypes(ct ...string) gzipFuncOpt {
	return func(opts *gzipOptions) {
		opts.contentTypes = append(opts.contentTypes, ct...)
	}
}

func MinLenght(len int) gzipFuncOpt {
	return func(opts *gzipOptions) {
		opts.MinLenght = len
	}
}

func CompressionLevel(level int) gzipFuncOpt {
	return func(opt *gzipOptions) {
		opt.CompressionLevel = level
	}

}

// GZIPBestCompression устанавливаем компрессию на уровень gzip.BestSpeed
func GZIPBestSpeed(opt *gzipOptions) {
	opt.CompressionLevel = gzip.BestSpeed
}

// GZIPBestCompression устанавливаем компрессию на уровень gzip.BestCompression
func GZIPBestCompression(opt *gzipOptions) {
	opt.CompressionLevel = gzip.BestCompression
}

func StatusCodes(codes ...int) gzipFuncOpt {
	return func(opt *gzipOptions) {
		opt.allowedStatusCodes = append(opt.allowedStatusCodes, codes...)
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

		// если мы ничего не писали в статус код, то подразумевается 200
		// и более того, если лишний раз пробовать писать этот статус, то он будет ругаться
		if cw.statusCode == 0 {
			cw.statusCode = 200
		}

		// либо мы прошли проверки на архивацию checked=true, либо архивания отменилась ignore==true
		if cw.opts.CheckStatusCode(cw.statusCode) && cw.opts.CheckContentTypes(cw.originalWriter.Header()) {
			cw.checked = true

		} else {
			// мы не будем использовать компрессию, устанавливаем флаг ignore и записываем байты в обернутый врайтер
			cw.ignore = true
			if cw.statusCode != 200 {
				// todo выделить это в какой-то отедльный блок, не нравится что всюду одно это
				cw.originalWriter.WriteHeader(cw.statusCode)
			}
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
	if cw.statusCode != 200 {
		cw.originalWriter.WriteHeader(cw.statusCode)
	}

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

func (cw *CompressedWriter) Close() {
	// Если порог минимального количество байт так и не пройдет, то просто сбрасываем все в
	// обернутый врайтер
	if cw.archiveWriter == nil {
		if cw.statusCode != 200 {
			cw.originalWriter.WriteHeader(cw.statusCode)
		}
		io.Copy(cw.originalWriter, &cw.fistBytesBuffer)
	} else {
		// а если архиватор создан, то его надо закрыть
		cw.archiveWriter.Close()
	}
	// это уже не часть интерфейса, но нужна для нашей мидвари, чтобы закрыть gzip.Writer
}

var _ http.ResponseWriter = (*CompressedWriter)(nil)

func GZIP(funcOpts ...gzipFuncOpt) func(http.Handler) http.Handler {
	// Создаем структуру настроек для GZIP,
	// после применения к ней опцефункций,
	// она будет передана в замыкание и будет уже работать с мидлварью
	opts := gzipOptions{
		CompressionLevel: gzip.BestCompression,
	}
	for _, f := range funcOpts {
		f(&opts)
	}

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
}

// encoded возвращает true, если запрос r закодирован
// кодировщиком encoder. Например:
//
// if encoded(r,"gzip"){...}
func encoded(r *http.Request, encoder string) bool {
	for _, v := range r.Header.Values("Content-Encoding") {
		if strings.Contains(v, encoder) {
			return true
		}
	}
	return false
}

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
