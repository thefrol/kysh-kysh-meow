package handlers

import (
	"context"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
)

type ForPing struct {
	Pinger func(context.Context) error
}

func (p ForPing) Ping(w http.ResponseWriter, r *http.Request) {
	err := p.Pinger(r.Context())
	if err != nil {
		httpio.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не удалось достучаться до базы данных: %v", err)
	}

	w.Write([]byte("0*0 все хорошо"))
}
