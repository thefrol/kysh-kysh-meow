package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/thefrol/kysh-kysh-meow/internal/slices"
)

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

// Apply !

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
// или сделать GZIP.Default
