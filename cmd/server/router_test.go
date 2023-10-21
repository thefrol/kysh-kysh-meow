package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		body     string
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
			name: "positive gauge #4",
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
			name:     "positive gauge #5",
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
			name:     "json value positive #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":100}`,
			},
		},
		{
			name: "json value positive #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/200"},
			method: http.MethodPost,
			route:  "/value",
			body:   `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":200}`,
			},
		},
		{
			name: "json value positive #3",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodPost,
			route:  "/value",
			body:   `{"type":"counter","id":"test1"}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"counter","id":"test1","delta":200}`,
			},
		},

		{
			name:     "json value negative #1",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodGet,
			route:    "/value",
			body:     `{"type":"gauge","id":"test1"}`,
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "text/plain",
				body:        "", // if == "" dont check
			},
		},
		{
			name:     "json value negative #2",
			prePosts: []string{"/update/gauge/test1/100"},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"gauge","id":"test2"}`,
			response: testResponse{
				code:        http.StatusNotFound,
				ContentType: "",
				body:        "",
			},
		},
		{
			name:     "json value negative #3 wrong type",
			prePosts: []string{},
			method:   http.MethodPost,
			route:    "/value",
			body:     `{"type":"unknown","id":"test2"}`,
			response: testResponse{
				code:        http.StatusBadRequest,
				ContentType: "",
				body:        "",
			},
		},

		{
			name: "json update positive #1",
			prePosts: []string{"/update/counter/test1/100",
				"/update/counter/test1/100"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"counter","id":"test1","delta":100}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"counter","id":"test1","delta":300}`,
			},
		},
		{
			name: "json update positive #2",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"gauge","id":"test1","value":100.1}`,
			response: testResponse{
				code:        http.StatusOK,
				ContentType: "application/json",
				body:        `{"type":"gauge","id":"test1","value":100.1}`,
			},
		},
		{
			name: "json update negative #1 unknown metric",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"unknown","id":"test1","gauge":100.1}`,
			response: testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "json update negative #2 crazy json",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `{"type":"gauge","id":"test1","gauge":100.1}`,
			response: testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "json update negative #3 not a json",
			prePosts: []string{"/update/gauge/test1/100",
				"/update/gauge/test1/101.2"},
			method: http.MethodPost,
			route:  "/update",
			body:   `lalalal`,
			response: testResponse{
				code: http.StatusBadRequest,
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

			resp, err := client.R().SetBody(tt.body).Execute(tt.method, server.URL+tt.route)
			require.NoErrorf(t, err, "error on final request %+v", err)

			assert.Equal(t, tt.response.code, resp.StatusCode())
			if tt.response.ContentType != "" {
				assert.Contains(t, resp.Header().Get("Content-Type"), tt.response.ContentType)
			}

			if tt.response.body != "" {
				assert.JSONEq(t, tt.response.body, string(resp.Body()))
			}

		})
	}

}
