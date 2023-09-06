package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// sendStats отправляет данные из хранилища store на сервер url,
// возвращает ошибку если что-то пошло не так
func sendStorageMetrics(store storage.Storager, url string) error {
	var errors []error
	for _, key := range store.ListCounters() {
		value, _ := store.Counter(key)
		err := doRequest(url, "counter", key, value) //#TODO counter to some const
		if err != nil {
			errors = append(errors, err)
		}
	}

	for _, key := range store.ListGauges() {
		value, _ := store.Gauge(key)
		err := doRequest(url, "gauge", key, value) //#TODO counter to some const
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

// doRequest создает запрос на сервер по нужному марштуру для обновления указанной метрики
func doRequest(host, metric, name string, value fmt.Stringer) error {
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
