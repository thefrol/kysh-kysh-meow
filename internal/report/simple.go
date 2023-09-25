package report

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// WithSimpleProtocol отправляет данные из хранилища store на сервер url,
// испрользуя простой протокол отправки, то есть ПОСТ запросами с пустыми телами.
// где запросы идут на server.addr/update/%counter_type%/%counter_name%/%value%
//
// Возвращает ошибку, если хотя бы один запрос не отправился. Так же сбрасывает
// запросы по середине, то есть если второй не отправился - остальные десять он даже
// и пытаться не будет.
//
// Deprecated
func WithSimpleProtocol(store storage.Storager, url string) error {
	var errors []error
	for _, key := range store.ListCounters() {
		value, _ := store.Counter(key)
		err := DoRequest(url, metrica.CounterName, key, value) //#TODO counter to some const
		if err != nil {
			errors = append(errors, err)
		}
	}

	for _, key := range store.ListGauges() {
		value, _ := store.Gauge(key)
		err := DoRequest(url, metrica.GaugeName, key, value) //#TODO counter to some const
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		s := ""
		for _, e := range errors {
			s += e.Error() + "\n"
		}
		//s = strings.TrimRight(s, ",")
		return fmt.Errorf(s)
	}
	return nil
}

// DoRequest создает запрос на сервер по нужному марштуру для обновления указанной метрики
//
// Deprecated
func DoRequest(host, metric, name string, value fmt.Stringer) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", host, metric, name, value)
	r, err := http.Post(url, "text/plain", nil)
	if err == nil {
		//прочитать тело и закрыть запрос, чтобы переиспользовать открытые tcp соединия
		_, err := io.Copy(io.Discard, r.Body)
		defer r.Body.Close()
		if err != nil {
			return err
		}

	}
	time.Sleep(200 * time.Millisecond) //защита от map concurrent write
	return err
}
