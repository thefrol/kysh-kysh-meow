package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func Test_dropCounter(t *testing.T) {
	type args struct {
		store       storage.Storager
		counterName string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		runsCount     int
		counterValues map[string]metrica.Counter //проверяет значение этого счетчика
		foundValues   map[string]bool            // проверяет наличие такого счетчика,
	}{
		{
			name:          "zero runs",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     0,
			counterValues: map[string]metrica.Counter{},
			foundValues:   map[string]bool{metricPollCount: false},
		},
		{
			name:          "one run",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     1,
			counterValues: map[string]metrica.Counter{metricPollCount: metrica.Counter(0)},
			foundValues:   nil,
		},
		{
			name:          "three runs",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     3,
			counterValues: map[string]metrica.Counter{metricPollCount: metrica.Counter(0)},
			foundValues:   nil,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			for i := 0; i < tt.runsCount; i++ {
				incrementCounter(tt.args.store, tt.args.counterName)
			}

			dropCounter(tt.args.store, tt.args.counterName)
			for k, expect := range tt.counterValues {
				real, ok := tt.args.store.Counter(k)
				assert.Truef(t, ok, "Not found counter %v", k)
				assert.Equal(t, expect, real)
			}
			for k, expect := range tt.foundValues {
				_, ok := tt.args.store.Counter(k)
				assert.Equalf(t, expect, ok, "Counter %v should be found(%v), but got %v", k, expect, ok)
			}

		})
	}
}

func Test_updateCounter(t *testing.T) {
	type args struct {
		store       storage.Storager
		counterName string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		runsCount     int
		counterValues map[string]metrica.Counter //проверяет значение этого счетчика
		foundValues   map[string]bool            // проверяет наличие такого счетчика,
	}{
		{
			name:          "zero runs",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     0,
			counterValues: map[string]metrica.Counter{},
			foundValues:   map[string]bool{metricPollCount: false},
		},
		{
			name:          "one run",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     1,
			counterValues: map[string]metrica.Counter{metricPollCount: metrica.Counter(1)},
			foundValues:   nil,
		},
		{
			name:          "three runs",
			args:          args{store: storage.New(), counterName: metricPollCount},
			wantErr:       false,
			runsCount:     3,
			counterValues: map[string]metrica.Counter{metricPollCount: metrica.Counter(3)},
			foundValues:   nil,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			for i := 0; i < tt.runsCount; i++ {
				incrementCounter(tt.args.store, tt.args.counterName)
			}
			for k, expect := range tt.counterValues {
				real, ok := tt.args.store.Counter(k)
				assert.Truef(t, ok, "Not found counter %v", k)
				assert.Equal(t, expect, real)
			}
			for k, expect := range tt.foundValues {
				_, ok := tt.args.store.Counter(k)
				assert.Equalf(t, expect, ok, "Counter %v should be found(%v), but got %v", k, expect, ok)
			}

		})
	}
}
