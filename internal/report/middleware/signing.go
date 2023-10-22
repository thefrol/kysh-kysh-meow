package middleware

import (
	"bytes"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/sign"
)

func Signing(key string) func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		br, ok := r.Body.(io.Reader)
		if !ok {
			log.Error().Msg("Не могу подписать данные, нужно чтобы в теле сообщение был ридер")
		}

		data := make([]byte, 500) // todo придумать как тут указать хороший слайс

		n, err := br.Read(data)
		if err != nil {
			return nil
		}

		s, err := sign.Bytes(data[:n], []byte(key))
		if err != nil {
			return err
		}
		r.Header.Set(sign.SignHeaderName, s)

		r.SetBody(bytes.NewBuffer(data[:n]))
		return nil
	}
}
