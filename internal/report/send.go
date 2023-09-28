package report

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/mailru/easyjson"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
)

var defaultClient = resty.New() // todo .SetJSONMarshaler(easyjson.Marshal())

// Send отправляет метрики из указанного хранилища store на сервер host.
// При возникновении ошибок будет стараться отправить как можно больше метрик,
// и продолжать работу, то есть, если первая метрика даст сбой, остальные двадцать он все же попытается отправить
// и вернет ошибку.
//
// При возникнвении ошибок возвращается только последняя
func Send(metricas []metrica.Metrica, url string) (lastErr error) {
	for _, m := range metricas {
		buf := new(bytes.Buffer)
		easyjson.MarshalToWriter(m, buf)
		// у нас существует очень важный контракт,
		// что тело сюда передается в формате io.Reader,
		// тогда могут работать разные мидлвари
		resp, err := defaultClient.R().SetBody(buf).Post(url) // todo в данный момент мы не используем тут easyjson

		if err != nil {
			ololog.Error().Str("location", "internal/report").Msgf("Не могу отправить сообщение с метрикой [%v]%v, по приничине %+v", m.MType, m.ID, err)
			lastErr = err
			continue
		}
		defer resp.RawBody().Close()
		ololog.Info().Msgf("Успешно отправлено %v %v", m.MType, m.ID)
	}
	return
}

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

// todo
//
// По-хорошему storage должен уметь вернуть мне все метрики в формате metrica.Metrica() и отказаться при этом от ListCounters()
//
// как будето на данный момент мы даже не можем представить какой-то список метрик, например которые были отправлены, или не были
//
// до сих пор не понимаю, что делать если счетчик PollCount не отправился, надо ли его сбрасывать, и вообще что делать...
// или у нас есть как бы текущая сессия? Мол складываем с исходным значением текущей сессии
// Но тогда это уже не REST
