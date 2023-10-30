package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/collector/internal/pollcount"
	"github.com/thefrol/kysh-kysh-meow/internal/compress"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
	"github.com/thefrol/kysh-kysh-meow/lib/retry"
	"github.com/thefrol/kysh-kysh-meow/lib/retry/fails"
)

var (
	ErrorRequestRejected = errors.New("запрос не принят, статус код не 200")
)

var defaultClient = resty.New()

// Send отправляет метрики из указанного хранилища store на сервер host.
// При возникновении ошибок будет стараться отправить как можно больше метрик,
// и продолжать работу, то есть, если первая метрика даст сбой, остальные двадцать он все же попытается отправить
// и вернет ошибку.
//
// При возникнвении ошибок возвращается только последняя
func Send(metricas []metrica.Metrica, url string) error {
	/*
		1. Замаршалить метрики в b
		2. Скомпрессировать полученные данные и заменить b,
			установить заголовок Content-Encoding
		3. Создать подпись на основе b, установить заголовок Sha256
		4. Попробовать отправить preparedRequest, если не получится, то ничего страшного
		5. Если получилось, обнуляем pollCounter
	*/
	preparedRequest := defaultClient.R()
	preparedRequest.Header.Set("Content-Type", "application/json")

	var b []byte // тут будет тело, которое в итоге запишем в сообщение

	b, err := json.Marshal(metricas)
	if err != nil {
		log.Error().Str("location", "internal/report").Msgf("Не могу замаршалить массив метрик по приничине %+v", err)
		return err
	}

	if len(b) >= int(CompressMinLength) {
		b, err = compress.Bytes(b, CompressLevel)
		if err != nil {
			return fmt.Errorf("ошибка компрессии: %w", err)
		}
		preparedRequest.Header.Set("Content-Encoding", "gzip")
	}

	if len(signingKey) != 0 {
		s, err := sign.Bytes(b, []byte(signingKey))
		if err != nil {
			return fmt.Errorf("ошибка подписывания: %w", err)
		}
		preparedRequest.Header.Set(sign.SignHeaderName, s)
		log.Info().Str("sign", s).Msg("Тело сообщения подписано")

		// мда, канеш цену за отсуствие мидлвари приходится платить
		// в таких не вполне очевидных ветвлениях
	}

	// подготавливаем запрос, в который теперь не будут вмешиваться мидвари
	preparedRequest.SetBody(b)
	var resp *resty.Response
	sendCall := func() error {
		var err error
		resp, err = preparedRequest.Post(url)
		return err
		// bug
		//
		// Есть небольшая проблема с этим resty:
		// Если будет ошибка распаковки сообщени после gzip,
		// то она звучит вот так "gzip: invalid header", а это
		// значит что мы воспримем это как ошибку отправки.
		// то есть мы даже не знаем получил ли метрики сервер,
		// а ошибка звучит как ошибка отправки
	}

	// запустим отправку с тремя попытками дополнительными
	err = retry.This(sendCall,
		retry.Attempts(3),
		retry.DelaySeconds(1, 3, 5, 7),
		retry.If(fails.OnDial), // черт, так красиво
		retry.OnRetry(func(n int, err error) {
			log.Info().Msgf("%v попытка отправить данные, ошибка: %v", n, err)
		}),
	)
	// todo в данный момент мы не используем тут easyjson

	if err != nil {
		log.Error().Str("location", "internal/report").Msgf("Не могу отправить сообщение c пачкой метрик по приничине %+v", err)
		return err
	}
	defer resp.RawBody().Close()

	log.Info().Str("location", "internal/report").Msgf("Метрики отправлены. Статус ответа %v, размер ответа %v", resp.StatusCode(), resp.Size())

	if resp.StatusCode() != http.StatusOK {
		log.Info().
			Str("location", "internal/report").
			Str("server_response", string(resp.Body())).
			Msg("Метрики отправлены, но не получены.")

		return fmt.Errorf("%w: сервер вернул %v: %v", ErrorRequestRejected, resp.StatusCode(), string(resp.Body()))
	}

	// Если сервер принял, то сбрасываем счетчик
	pollcount.Drop()

	return nil
}

// CompressLevel устанавливает минимальное число байт в теле сообения, после которого
// начинается комапрессия
var CompressMinLength uint = 100

// CompressLevel устанавливает профиль сжатия gzip, по умолчанию равен
// gzip.BestCompression
var CompressLevel int = compress.BestCompression

// signingKey содержит ключ подписывания
var signingKey []byte

// SetSigningKey устанавливает ключ подписывания для отправляемых запросов
func SetSigningKey(key string) {
	signingKey = []byte(key)
}
