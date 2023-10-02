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
	"HeapReleased",
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
	tests := []struct {
		name           string
		wantErr        bool
		memValuesCount int      // if <0 not checking
		fieldsFound    []string //какие поля мы должны содержаться в memStore, можно неполно
	}{
		{
			name:           "all metrics in place",
			wantErr:        false,
			memValuesCount: -1,
			fieldsFound:    metricsMem,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			st := Fetch()
			if tt.memValuesCount >= 0 {
				assert.Equal(t, tt.memValuesCount, len(st.ToTransport()))
			}
			for _, v := range tt.fieldsFound {
				assert.Truef(t, findMetric(st, "gauge", randomValueName) || findMetric(st, "gauge", randomValueName), "Not found metric %v", v)
			}

		})
	}
}

// Test_fetchAdditionalStats проверяет, что случайная величина так же хорошо сохраняется в хранилище
func Test_fetchAdditionalStats(t *testing.T) {

	tests := []struct {
		name        string
		wantErr     bool
		fieldsFound []string //какие поля мы должны содержаться в memStore, можно неполно
	}{
		{
			name:        "all metrics in place",
			wantErr:     false,
			fieldsFound: []string{randomValueName},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			st := Fetch()
			for _, v := range tt.fieldsFound {
				assert.Truef(t, findMetric(st, "gauge", randomValueName), "Not found metric %v", v)
			}

		})
	}
}

func CountValues(s storage.Storager) int {
	return len(s.ListCounters()) + len(s.ListGauges())
}

func findMetric(st Stats, mtype string, name string) bool {
	for _, v := range st.ToTransport() {
		if v.ID == name && v.MType == mtype {
			return true
		}
	}
	return false
}
