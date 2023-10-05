package report

import (
	"bytes"

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
		_, err := easyjson.MarshalToWriter(m, buf)
		if err != nil {
			ololog.Error().Str("location", "internal/report").Msgf("Не могу замаршалить [%v]%v, по приничине %+v", m.MType, m.ID, err)
			lastErr = err
			continue
		}
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

	// сбрасываем счетчик PollCount
	dropPollCount() // todo вот этот сброс надо проверять
	return
}

// UseBeforeRequest встраивает мидлварь в цепочку отправки сообщений. Все обработчики получают доступ
// к рести клиенту и текущему подготавливаемому запросу. Таким образом можно сделать дополнительное поггирование,
// или сжатие
//
// пример: report.UseBeforeRequest(GZIP)
func UseBeforeRequest(middlewares ...func(c *resty.Client, r *resty.Request) error) {
	for _, m := range middlewares {
		defaultClient.OnBeforeRequest(m)
	}
}
