// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func Test_saveMemStats(t *testing.T) {
	type args struct {
		store   storage.Storager
		exclude []string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		memValuesCount int      // if <0 not checking
		fieldsFound    []string //какие поля мы должны содержаться в memStore, можно неполно
	}{
		{
			name:           "all metrics in place",
			args:           args{store: storage.New()},
			wantErr:        false,
			memValuesCount: -1,
			fieldsFound:    metricsMem,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := saveMemStats(tt.args.store, tt.args.exclude); (err != nil) != tt.wantErr {
				t.Errorf("saveMemStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.memValuesCount >= 0 {
				assert.Equal(t, tt.memValuesCount, CountValues(tt.args.store))
			}
			for _, v := range tt.fieldsFound {
				_, gaugeFound := tt.args.store.Gauge(v)
				_, counterFound := tt.args.store.Counter(v)
				assert.Truef(t, gaugeFound || counterFound, "Not found metric %v", v)
			}

		})
	}
}

func Test_saveAdditionalStats(t *testing.T) {
	type args struct {
		store storage.Storager
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		fieldsFound []string //какие поля мы должны содержаться в memStore, можно неполно
	}{
		{
			name:        "all metrics in place",
			args:        args{store: storage.New()},
			wantErr:     false,
			fieldsFound: []string{metricRandomValue},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if err := saveAdditionalStats(tt.args.store); (err != nil) != tt.wantErr {
				t.Errorf("saveMemStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			for _, v := range tt.fieldsFound {
				_, gaugeFound := tt.args.store.Gauge(v)
				_, counterFound := tt.args.store.Counter(v)
				assert.Truef(t, gaugeFound || counterFound, "Not found metric %v", v)
			}

		})
	}
}

func Test_updateCounter(t *testing.T) {
	type args struct {
		store storage.Storager
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		runsCount     int
		counterValues map[string]storage.Counter //проверяет значение этого счетчика
		foundValues   map[string]bool            // проверяет наличие такого счетчика,
	}{
		{
			name:          "zero runs",
			args:          args{store: storage.New()},
			wantErr:       false,
			runsCount:     0,
			counterValues: map[string]storage.Counter{},
			foundValues:   map[string]bool{metricPollCount: false},
		},
		{
			name:          "one run",
			args:          args{store: storage.New()},
			wantErr:       false,
			runsCount:     1,
			counterValues: map[string]storage.Counter{metricPollCount: storage.Counter(1)},
			foundValues:   nil,
		},
		{
			name:          "three runs",
			args:          args{store: storage.New()},
			wantErr:       false,
			runsCount:     3,
			counterValues: map[string]storage.Counter{metricPollCount: storage.Counter(3)},
			foundValues:   nil,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			for i := 0; i < tt.runsCount; i++ {
				if err := updateCounter(tt.args.store); (err != nil) != tt.wantErr {
					t.Errorf("saveMemStats() error = %v, wantErr %v", err, tt.wantErr)
				}
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

func CountValues(s storage.Storager) int {
	return len(s.ListCounters()) + len(s.ListGauges())
}
