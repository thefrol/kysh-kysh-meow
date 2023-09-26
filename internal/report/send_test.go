package report_test

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/report"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name     string
		metricas []metrica.Metrica
		wantErr  bool
	}{
		{
			name: "counters",
			metricas: []metrica.Metrica{
				{ID: "test1", Delta: wrapInt64(10), MType: "counter"},
				{ID: "Test1", Delta: wrapInt64(9), MType: "counter"}},
			wantErr: false,
		},
		{
			name: "gauges",
			metricas: []metrica.Metrica{
				{ID: "test1", Value: wrapFloat64(10), MType: "gauge"},
				{ID: "Test1", Value: wrapFloat64(9), MType: "gauge"}},
			wantErr: false,
		},
		{
			name: "unknown metricas",
			metricas: []metrica.Metrica{
				{ID: "test1", Value: wrapFloat64(10), MType: "unknown"},
				{ID: "Test22", Value: wrapFloat64(9), MType: "unknown"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testHandler{}
			s := httptest.NewServer(&h)

			err := report.Send(tt.metricas, s.URL)
			if tt.wantErr {
				assert.Error(t, err, "Должна быть ошибка")
			}

			var sended []metrica.Metrica
			for _, r := range h.requests {
				m := new(metrica.Metrica)
				body, err := r.GetBody()
				require.NoError(t, err, "Тест сервер не смог прочитать тело полученного запроса")
				err = easyjson.UnmarshalFromReader(body, m)
				require.NoError(t, err, "Не возможно размаршалить джейсон из отправленного джейсона")
				defer r.Body.Close()
				sended = append(sended, *m)
			}

			eq := reflect.DeepEqual(tt.metricas, sended)
			assert.True(t, eq, "Ожидаемые к отправке данные не совпадают с полученными сервером")
		})
	}
}

func wrapInt64(v int) (ref *int64) {
	ref = new(int64)
	*ref = int64(v)
	return
}

func wrapFloat64(v int) (ref *float64) {
	ref = new(float64)
	*ref = float64(v)
	return
}
