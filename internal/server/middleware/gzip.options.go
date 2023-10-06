package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/thefrol/kysh-kysh-meow/lib/slices"
)

type gzipOptions struct {
	MinLenght          int
	contentTypes       []string
	CompressionLevel   int
	allowedStatusCodes []int
}

// Apply применяеn опцефункцию к настройкам GZIP
func (o *gzipOptions) Apply(opts ...gzipFuncOpt) {
	for _, f := range opts {
		f(o)
	}
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

// todo По-хорошему, у нас должны быть ещё какие-то удалятели значений, типа BlockStatusCode(200), типа чтобы можно было удалить из установленных автоматически или типа того

// GZIPDefault создает типичные настройки для сжатия:
//
//	Уровень сжатия: лучшая компрессия
//	Заголовки: text/html, text/plain, text/xml, text/css, application/json, application/javascript
//	Статус коды: 200
//	Минимальная длинна 50 байт
func GZIPDefault(opt *gzipOptions) {
	opt.Apply(
		GZIPBestCompression,
		StatusCodes(http.StatusOK),
		ContentTypes(
			"text/plain",
			"text/html",
			"text/css",
			"text/xml",
			"application/json",
			"application/javascript"),
		MinLenght(50))
}