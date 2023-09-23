// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/slices"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
	"github.com/thefrol/kysh-kysh-meow/internal/structs"
)

const (
	metricPollCount   = "PollCount"
	metricRandomValue = "RandomValue"
)

// fetchStats собирает метрики мамяти и сохраняет их в хранилище, исключая ненужные exclude
func saveMemStats(store storage.Storager, exclude []string) error {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	stats, err := structs.FieldsFloat(m, exclude)
	if err != nil {
		return err
	}
	for key, value := range stats {
		store.SetGauge(key, metrica.Gauge(value)) //#TODO SetGauges()
	}
	return nil
}

// saveAdditionalStats сохраняет дополнительные метрики в указанное хранилище:
// 	RandomValue: тип Gauge, содержит в себе случайное значение
func saveAdditionalStats(store storage.Storager) error {
	store.SetGauge(metricRandomValue, metrica.Gauge(randomFloat64()))
	return nil
}

// parseMetrics возвращает метрики, не найденные в выдаче runtime(lost) и список несохраняемых метрик(excluded).
// В случае ошибки вернет ошибку третим параметром
func parseMetrics() (lost []string, exclude []string, err error) {
	m := runtime.MemStats{} // по хорошему тут бы получать какой-то пустой элемент, и чтобы getFields работала только с типом!
	fields, err := structs.FieldNames(m)
	if err != nil {
		return nil, nil, fmt.Errorf("can't retrieve fields from MemStats")
	}
	lost = slices.Difference[string](metricsMem, fields)
	exclude = slices.Difference[string](fields, metricsMem)
	return lost, exclude, nil
}

// randomFloat64 возвращает случайное число типа float64
func randomFloat64() float64 {
	s := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(s)
	return r.Float64()
}
