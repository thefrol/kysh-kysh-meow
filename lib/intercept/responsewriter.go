package intercept

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Необходимо закрыть
type WriteInterceptor struct {
	statusCode int
	origWriter http.ResponseWriter
	buf        *bytes.Buffer
}

func New(w http.ResponseWriter, data []byte) WriteInterceptor {
	return WriteInterceptor{
		origWriter: w,
		buf:        bytes.NewBuffer(data[:0]), // обнуляем массив, иначе пишет в конец, а читать будет с начала, и в выход выйдет мусор, что бы уже в буфере
	}
}

func (w *WriteInterceptor) WriteHeader(code int) {
	w.statusCode = code
}

func (w *WriteInterceptor) Header() http.Header {
	return w.origWriter.Header()
}

func (w WriteInterceptor) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

// Close говорит о том, что сообщение готовится к отправке, значит
// можно установить нужные хедеры и отправлять
func (w WriteInterceptor) Close() {
	if w.statusCode != 0 {
		w.origWriter.WriteHeader(w.statusCode)
	}

	_, err := io.Copy(w.origWriter, w.buf)
	if err != nil {
		log.Error().Msgf("copy a response to originalWriter: %v", err)
	}
}

func (w WriteInterceptor) Buf() *bytes.Buffer {
	return w.buf
}
