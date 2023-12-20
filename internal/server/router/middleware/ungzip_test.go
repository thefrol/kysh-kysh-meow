package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUncompress(t *testing.T) {
	type request struct {
		contentEncoding string
		encodeGZIP      bool
		body            string
	}
	type response struct {
		body string // when =="" not testing
		code int
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		{
			name:     "basic",
			request:  request{body: "anything", contentEncoding: "gzip", encodeGZIP: true},
			response: response{body: "anything", code: http.StatusOK},
		},
		{
			name:     "pass on another encoder",
			request:  request{body: "anything", contentEncoding: "deflate"},
			response: response{body: "anything", code: http.StatusOK},
		},
		{
			name:     "badly encoded body",
			request:  request{body: "blabla-zip-not-encoded", contentEncoding: "gzip"},
			response: response{code: http.StatusBadRequest},
		},
	}
	r := chi.NewRouter()
	r.Use(UnGZIP)
	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err, "Cant read body in inner handler of test")
		w.Write(buf)
	}))
	client := resty.New()
	serv := httptest.NewServer(r)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body string
			if tt.request.encodeGZIP {
				b := new(bytes.Buffer) // todo можно ли обойтись без b
				gzw := gzip.NewWriter(b)
				gzw.Write([]byte(tt.request.body))
				gzw.Close()
				body = b.String()
			} else {
				body = tt.request.body
			}

			resp, err := client.R().
				SetHeader("Content-Encoding", tt.request.contentEncoding).
				SetBody(body).
				Execute(http.MethodPost, serv.URL)
			require.NoError(t, err, "An error on request")
			assert.Equal(t, tt.response.code, resp.StatusCode())
			//skip if body is empty
			if tt.response.body != "" {
				assert.Equal(t, tt.response.body, string(resp.Body()))
			}

		})
	}
}
