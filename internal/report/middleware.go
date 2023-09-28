package report

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

// UseBeforeRequest встраиваем мидлварь в цепочку отправки сообщений. Все обработчики получают доступ
// к рести клиенту и текущему подготавливаемому запросу. Таким образом можно сделать дополнительное поггирование,
// или сжатие
//
// пример: report.UseBeforeRequest(GZIP)
func UseBeforeRequest(middlewares ...func(c *resty.Client, r *resty.Request) error) {
	for _, m := range middlewares {
		defaultClient.OnBeforeRequest(m)
	}
}

func ApplyGZIP(minLenght int, level int) func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		// проверяем, что контент уже не закодирован каким-нибудт другим мидлварью или кодом
		if r.Header.Values("Content-Encoding") != nil {
			ololog.Warn().Str("location", "agent/middleware/gzip").Fields(r.Header).Msg("Запрос на сервер уже сжат")
			return nil
		}

		v, ok := r.Body.(io.Reader)
		if !ok {
			ololog.Warn().Str("location", "agent/middleware/gzip").Msg("Тело сообщение передано не в формате io.Reader, а значит мы не можем его заархивировать. Вообще мы хотим передавать тело сообщения именно ридером")
			return nil
		}
		// todo осталось проверить, что он минимальная длинна достигнута, попробовать прочитать первые сколько то байт
		// на данный момент сообщения длинной в 20-30 байт стали длинной 50-70 байт
		b := new(bytes.Buffer)
		gz, err := gzip.NewWriterLevel(b, level)
		if err != nil {
			return fmt.Errorf("cant create compressor")
		}
		_, err = io.Copy(gz, v)
		if err != nil {
			return fmt.Errorf("cant write to gzip writer")
		}
		r.Header.Add("Content-Encoding", "gzip")
		r.SetBody(b)
		gz.Close()

		return nil
	}
}
