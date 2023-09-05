package main

import (
	"fmt"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// sendStats отправляет данные из хранилища store на сервер url,
// возвращает ошибку если что-то пошло не так
func sendStorageMetrics(store storage.Storager, url string) error {
	var errors []error
	for _, key := range store.ListCounters() {
		value, _ := store.Counter(key)
		err := sendMetric(url, "counter", key, value) //#TODO counter to some const
		if err != nil {
			errors = append(errors, err)
		}
	}

	for _, key := range store.ListGauges() {
		value, _ := store.Gauge(key)
		err := sendMetric(url, "gauge", key, value) //#TODO counter to some const
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

func sendMetric(host, metric, name string, value fmt.Stringer) error {
	url := fmt.Sprintf("%s/update/%s/%s/%s", host, metric, name, value)
	_, err := http.Post(url, "text/plain", nil)
	return err
}
