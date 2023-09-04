// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"runtime"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// fetchStats собирает метрики мамяти и сохраняет их в хранилище, исключая ненужные exclude
func saveMemStats(store storage.Storager, exclude []string) error {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	stats, err := getFieldsFloat(m, exclude)
	if err != nil {
		return err
	}
	for key, value := range stats {
		store.SetGauge(key, storage.Gauge(value)) //#TODO SetGauges()
	}
	return nil
}
