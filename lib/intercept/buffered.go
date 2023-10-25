package intercept

import (
	"io"
	"net/http"
)

// Buffered перехватывает все данные, которые должны быть
// записаны в респонс врайтер, и хранит в памяти, а записывает
// в момент вызова Flush(). Обязательно в конце вызывать Flush().
//
// Применяется, когда мы хотим получить тело сообщения, и после
// провести с ним некие процедуры, например подписать, или архивировать.
// Но если во врайтер уже что-то записано, то мы не сможем поменять заголовки или
// или выход, поэтому мы перехватываем записанные данные, и пишем их в конце
//
// Заголовки при этом не перехватываются, а записываются в оригинальный
// врайтер сразу
//
// поведение defer w.Flush() не проверено
type Buffered struct {
	http.ResponseWriter
	buf  io.ReadWriter
	code int
}

// WithBuffer создает буферную оболочку для w, где в оригинальный
// врайтер иничего не будет записано до вызова Flush(),
// было бы правильно называть этот вызов Close()
func WithBuffer(w http.ResponseWriter, buf io.ReadWriter) *Buffered {
	return &Buffered{
		ResponseWriter: w,
		buf:            buf,
	}
}

func (w *Buffered) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *Buffered) WriteHeader(code int) {
	w.code = code
}

// Flush записывает перехваченные данные, такие как код и тело ответа
// в оригинальный врайтер
func (w Buffered) Flush() error {
	if w.code != 0 {
		w.WriteHeader(w.code)

	}
	_, err := io.Copy(w.ResponseWriter, w.buf)
	return err
}

func (w Buffered) StatusCode() int {
	return w.code
}
