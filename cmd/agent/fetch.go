// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
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

func saveAdditionalStats(store storage.Storager) error {
	store.SetGauge(metricRandomValue, metrica.Gauge(randomFloat64()))
	return nil
}

func updateCounter(store storage.Storager) error {
	count, _ := store.Counter(metricPollCount)
	store.SetCounter(metricPollCount, count+metrica.Counter(1))
	return nil
}

func randomFloat64() float64 {
	s := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(s)
	return r.Float64()
}
