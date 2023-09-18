// fetchStats собирает метрики из памяти, и выдает их удобной мапой
package stats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// Список сохраняемых метрик из пакета runtime
var metricsMem = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

func Test_fetchMemStats(t *testing.T) {
	type args struct {
		store storage.Storager
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
			FetchMemStats(tt.args.store)
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

// Test_fetchAdditionalStats проверяет, что случайная величина так же хорошо сохраняется в хранилище
func Test_fetchAdditionalStats(t *testing.T) {
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
			fieldsFound: []string{randomValueName},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			FetchMemStats(tt.args.store)
			for _, v := range tt.fieldsFound {
				_, gaugeFound := tt.args.store.Gauge(v)
				_, counterFound := tt.args.store.Counter(v)
				assert.Truef(t, gaugeFound || counterFound, "Not found metric %v", v)
			}

		})
	}
}

func CountValues(s storage.Storager) int {
	return len(s.ListCounters()) + len(s.ListGauges())
}
