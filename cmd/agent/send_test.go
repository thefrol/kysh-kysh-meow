package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/slices"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func Test_sendMetric(t *testing.T) {
	type args struct {
		//host   string
		metric string
		name   string
		value  fmt.Stringer
	}
	tests := []struct {
		name          string
		serverHost    string
		args          args
		routesUsed    []string
		requestsCount int //если меньше одного, то не проверяется
		wantErr       bool
	}{
		{
			name:          "positive #1",
			args:          args{metric: "counter", name: "test1", value: stringerWrap{"1234"}},
			routesUsed:    []string{"/update/counter/test1/1234"},
			requestsCount: 1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := testServer{}
			s := httptest.NewServer(&ts)
			defer s.Close()

			if err := doRequest(s.URL, tt.args.metric, tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.requestsCount, ts.NumRequests(), "Wrong requests count")
			//assert.Contains(t, server.routesUsed(), tt.routesUsed, "Not all required routes used") // не срабатывает, пишешь {"1"} not containg {"1"}
			assert.True(t, slices.СontainsSlice[string](ts.routesUsed(), tt.routesUsed))
		})
	}
}

func Test_sendStorageMetrics(t *testing.T) {
	type args struct {
		counters map[string]metrica.Counter
		gauges   map[string]metrica.Gauge
	}
	tests := []struct {
		name          string
		serverHost    string
		args          args
		routesUsed    []string
		requestsCount int //если меньше одного, то не проверяется
		wantErr       bool
	}{
		{
			name: "positive #1 counter",
			args: args{
				counters: map[string]metrica.Counter{
					"test1": metrica.Counter(22)}},
			routesUsed:    []string{"/update/counter/test1/.*"},
			requestsCount: 1,
			wantErr:       false,
		},
		{
			name:       "positive #2 gauge",
			serverHost: newHost(),
			args: args{
				gauges: map[string]metrica.Gauge{
					"test1": metrica.Gauge(22.1)}},
			routesUsed:    []string{"/update/gauge/test1/.*"},
			requestsCount: 1,
			wantErr:       false,
		},
		//мне нужен тест на неправильный адрес! для тестирования ошибок или типа того. МОгу ли я это протестировать правда? или просто гонюсь за покрытием?
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.New()

			// устанавливаем значения счетчиков
			for k, v := range tt.args.counters {
				store.SetCounter(k, v)
			}
			for k, v := range tt.args.gauges {
				store.SetGauge(k, v)
			}

			// Запускаем сервак
			ts := testServer{}
			s := httptest.NewServer(&ts)
			defer s.Close()

			if err := sendStorageMetrics(store, s.URL); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.requestsCount, ts.NumRequests(), "Wrong requests count")
			//assert.Contains(t, server.routesUsed(), tt.routesUsed, "Not all required routes used") // не срабатывает, пишешь {"1"} not containg {"1"}
			for _, v := range tt.routesUsed {
				found, err := ts.containsRouteRegexp(v)
				assert.NoError(t, err)
				assert.Truef(t, found, "routesUsed %v must contain %v", ts.routesUsed(), v)
			}
		})
	}
}

// Хендлеры сервера вынести в отдельный пакет и тестировать их тут #TODO

// testServer простая структура, которая запоминает по каким адресам к ней делали запросы
// используется в тестировании
type testServer struct {
	requests []*http.Request
}

// ServerHTTP отвечает на запросы к серверы, исполяет интерфейс http.Handle
func (server *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got a request")
	server.requests = append(server.requests, r)
}

// routesUsed возвращает марштуры по которым проходили запросы
func (server testServer) routesUsed() (routes []string) {
	for _, r := range server.requests {
		routes = append(routes, r.URL.Path)
	}
	return
}

func (server testServer) containsRouteRegexp(pattern string) (bool, error) {
	for _, r := range server.requests {
		found, err := regexp.MatchString(pattern, r.URL.Path)
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}
	return false, nil
}

// возвращает количество полученных запросов
func (server testServer) NumRequests() int {
	return len(server.requests)
}

// stringerWrap позволяет использовать интерфейс стрингер в полях структуры. Позволяя запихивать туда люую переменную, которую можно обратить в строку
type stringerWrap struct {
	text string
}

func (s stringerWrap) String() string {
	return s.text
}

var freeport int = 8081

func newHost() (url string) {
	freeport++
	url = lastHost()
	return
}

func lastHost() string {
	return fmt.Sprintf("localhost:%v", freeport)
}
