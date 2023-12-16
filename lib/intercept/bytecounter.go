package intercept

import "net/http"

type BytesCounter struct {
	http.ResponseWriter
	n    int
	code int
}

// Write воплощает интерфейс http.ResponseWriter;
// При каждом вызове, увеличиваем счетчик записанных байт,
// при этом данные пишутся в оригинальный врайтер
func (w *BytesCounter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.n += n
	return n, err
}

func (w *BytesCounter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w BytesCounter) BytesWritten() int {
	return w.n
}

func (w BytesCounter) StatusCode() int {
	return w.code
}

func WithBytesCounter(w http.ResponseWriter) *BytesCounter {
	return &BytesCounter{ResponseWriter: w, code: 200}
}
