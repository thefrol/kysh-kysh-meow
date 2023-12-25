package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storagev2/mem"
)

// BenchmarkHandlers проверяет как работают хендлеры.
//
// В качестве хранилища используется хранилище в памяти
// чтобы в бенчмарки не попадали проблемы с БД
func BenchmarkHandlers(b *testing.B) {

	// айдишники счетчиков, которые будут использоваться
	const (
		metricID = "testName" // используется и для гаужей и для счетчиков
	)

	s := mem.MemStore{
		Counters: make(mem.IntMap, 50),
		Gauges:   make(mem.FloatMap, 50),
	}

	s.CounterIncrement(context.Background(), metricID, 10)
	s.GaugeUpdate(context.Background(), metricID, 10)

	m := manager.Registry{
		Counters: &s,
		Gauges:   &s,
	}

	j := metricas.Manager{
		Registry: m,
	}

	api := router.API{
		Registry: m,
		Manager:  j,
	}

	h := api.MeowRouter()

	v := int64(10)
	route := "/value/"
	rec := httptest.NewRecorder()

	me := metrica.Metrica{
		MType: "counter",
		ID:    metricID,
		Delta: &v,
	}
	js, err := json.Marshal(&me)
	if err != nil {
		b.Error("cant marshall metric")
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
	h.ServeHTTP(rec, req)

	// уберем лишний вывод
	log.Logger = log.Level(zerolog.WarnLevel)

	b.Run("getting counters with query", func(b *testing.B) {
		route := fmt.Sprintf("/value/counter/%v", metricID)
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.ServeHTTP(rec, req)
		}
	})

	b.Run("getting gauges with query", func(b *testing.B) {
		route := fmt.Sprintf("/value/gauge/%v", metricID)
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.ServeHTTP(rec, req)
		}
	})

	b.Run("updating gauges with query", func(b *testing.B) {
		v := 140
		route := fmt.Sprintf("/update/gauge/%v/%v", metricID, v)
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.ServeHTTP(rec, req)
		}
	})

	b.Run("updating counters with query", func(b *testing.B) {
		v := 140
		route := fmt.Sprintf("/update/counter/%v/%v", metricID, v)
		req := httptest.NewRequest(http.MethodGet, route, nil)
		rec := httptest.NewRecorder()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.ServeHTTP(rec, req)
		}
	})

	b.Run("getting counter with json", func(b *testing.B) {
		route := "/value"
		rec := httptest.NewRecorder()

		m := metrica.Metrica{
			MType: "counter",
			ID:    metricID,
		}
		js, err := json.Marshal(&m)
		if err != nil {
			b.Error("cant marshall metric")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
			b.StartTimer()

			h.ServeHTTP(rec, req)
		}
	})

	b.Run("updating counter with json", func(b *testing.B) {
		v := int64(10)
		route := "/update"
		rec := httptest.NewRecorder()

		m := metrica.Metrica{
			MType: "counter",
			ID:    metricID,
			Delta: &v,
		}
		js, err := json.Marshal(&m)
		if err != nil {
			b.Error("cant marshall metric")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
			b.StartTimer()

			h.ServeHTTP(rec, req)
		}
	})

	b.Run("getting gauge with json", func(b *testing.B) {
		route := "/value"
		rec := httptest.NewRecorder()

		m := metrica.Metrica{
			MType: "gauge",
			ID:    metricID,
		}
		js, err := json.Marshal(&m)
		if err != nil {
			b.Error("cant marshall metric")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
			b.StartTimer()

			h.ServeHTTP(rec, req)
		}
	})

	b.Run("updating gauge with json", func(b *testing.B) {
		v := float64(10)
		route := "/update"
		rec := httptest.NewRecorder()

		m := metrica.Metrica{
			MType: "gauge",
			ID:    metricID,
			Value: &v,
		}
		js, err := json.Marshal(&m)
		if err != nil {
			b.Error("cant marshall metric")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
			b.StartTimer()

			h.ServeHTTP(rec, req)
		}
	})

	b.Run("batch update with json", func(b *testing.B) {
		v := float64(10)
		route := "/updates"
		rec := httptest.NewRecorder()

		m := metrica.Metrica{
			MType: "gauge",
			ID:    metricID,
			Value: &v,
		}

		ba := metrica.Metricas{m}

		js, err := json.Marshal(&ba)
		if err != nil {
			b.Error("cant marshall metric")
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			req := httptest.NewRequest(http.MethodPost, route, bytes.NewBuffer(js))
			b.StartTimer()

			h.ServeHTTP(rec, req)
		}
	})

}
