package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

func Test_MeowRouter(t *testing.T) {

	type testResponse struct {
		code        int
		ContentType string
		body        string // сейчас не используется
	}

	tests := []struct {
		name     string
		prePosts []string // urls to call before the main test
		method   string
		response testResponse
		route    string
	}{
		{
			name:     "positive counter #1",
			prePosts: []string{"/update/counter/test1/100"},
			method:   http.MethodGet,
			route:    "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name: "positive counter #2",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodGet,
			route:  "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "200",
			},
		},
		{
			name: "positive counter #3",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100.1"}, // float значение сервер не пример
			method: http.MethodGet,
			route:  "/value/counter/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name:     "positive gauge #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
		{
			name: "positive gauge #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1"},
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100.1",
			},
		},
		{
			name: "positive gauge #3",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1",
				"/update/gauge/test1/none"}, //это значение не примет
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100.1",
			},
		},
		{
			name: "positive gauge #3",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/100.1",
				"/update/gauge/test1/2e0"},
			method: http.MethodGet,
			route:  "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "2",
			},
		},
		{
			name:     "positive gauge #4",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value/gauge/test1",
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "text/plain",
				body:        "100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store = storage.New() //обнуляем хранилище
			server := httptest.NewServer(MeowRouter())
			defer server.Close()
			client := resty.New()
			for _, u := range tt.prePosts {
				_, err := client.R().SetHeader("Content-Type", "text/plain").Post(server.URL + u)
				assert.NoError(t, err, "error on preparing data for final request")
			}

			resp, err := client.R().Execute(tt.method, server.URL+tt.route)
			assert.NoError(t, err, "error on final request")

			assert.Equal(t, tt.response.code, resp.StatusCode())
			assert.Contains(t, resp.Header().Get("Content-Type"), tt.response.ContentType)
			assert.Equal(t, tt.response.body, string(resp.Body()))

		})
	}

}
