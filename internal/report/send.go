package report

import (
	"github.com/go-resty/resty/v2"
	"github.com/thefrol/kysh-kysh-meow/internal/ololog"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var defaultClient = resty.New() // todo .SetJSONMarshaler(easyjson.Marshal())

// Send отправляет метрики из указанного хранилища store на сервер host.
// При возникновении ошибок будет стараться отправить как можно больше метрик,
// и продолжать работу, то есть, если первая метрика даст сбой, остальные двадцать он все же попытается отправить
// и вернет ошибку.
func Send(store storage.Storager, url string) (last_err error) {
	for _, m := range store.Metricas() {
		resp, err := defaultClient.R().SetBody(m).Post(url) // todo в данный момент мы не используем тут easyjson
		if err != nil {
			ololog.Error().Str("location", "internal/report").Msgf("Не могу отправить сообщение с метрикой [%v]%v, по приничине %+v", m.MType, m.ID, err)
			last_err = err
			continue
		}
		defer resp.RawBody().Close()
		ololog.Info().Msgf("Успешно отправлено %v %v", m.MType, m.ID)
	}
	return
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
