package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

var (
	ErrorsServerError = errors.New("ошибка сервера")
)

var defaultClient = resty.New() // todo .SetJSONMarshaler(easyjson.Marshal())

// Send отправляет метрики из указанного хранилища store на сервер host.
// При возникновении ошибок будет стараться отправить как можно больше метрик,
// и продолжать работу, то есть, если первая метрика даст сбой, остальные двадцать он все же попытается отправить
// и вернет ошибку.
//
// При возникнвении ошибок возвращается только последняя
func Send(metricas []metrica.Metrica, url string) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(metricas)
	if err != nil {
		log.Error().Str("location", "internal/report").Msgf("Не могу замаршалить массив метрик по приничине %+v", err)
		return err
	}
	// у нас существует очень важный контракт,
	// что тело сюда передается в формате io.Reader,
	// тогда могут работать разные мидлвари
	resp, err := defaultClient.R().SetBody(buf).Post(url) // todo в данный момент мы не используем тут easyjson

	if err != nil {
		log.Error().Str("location", "internal/report").Msgf("Не могу отправить сообщение c пачкой метрик по приничине %+v", err)
		return err
	}
	defer resp.RawBody().Close()

	log.Info().Str("location", "internal/report").Msgf("Метрики отправлены. Статус ответа %v, размер %v", resp.StatusCode(), resp.Size())

	if resp.StatusCode() != http.StatusOK {
		log.Info().
			Str("location", "internal/report").
			Str("server_response", string(resp.Body())).
			Msg("Метрики отправлены, но не получены.")

		return fmt.Errorf("%w: сервер вернул %v: %v", ErrorsServerError, resp.StatusCode(), resp.Body())
	}

	// Если сервер принял, то сбрасываем счетчик
	dropPollCount()

	return nil
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
