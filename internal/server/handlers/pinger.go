package handlers

import (
	"context"
	"net/http"

	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

type ForPing struct {
	Pinger func(context.Context) error
}

func (p ForPing) Ping(w http.ResponseWriter, r *http.Request) {
	err := p.Pinger(r.Context())
	if err != nil {
		api.HTTPErrorWithLogging(w,
			http.StatusInternalServerError,
			"Не удалось достучаться до базы данных: %v", err)
	}

	w.Write([]byte("0*0 все хорошо"))
}
