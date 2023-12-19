package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
	"github.com/thefrol/kysh-kysh-meow/internal/server/httpio"
)

func (a *ForQuery) GetCounter(w http.ResponseWriter, r *http.Request) {
	var (
		id = chi.URLParam(r, "id")
	)

	// ментор интересно, нрмально если мы под каждый хендлер будем делать свой класс и у него уже будет
	// логгер на готове и все такое
	logger := log.With().
		Str("handler", "GetCounter()").
		Str("id", id).
		Logger()

	v, err := a.Registry.Counter(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorMetricNotFound) {
			logger.Error().Err(err).Send() // httpio.To(log).To(w).NotFound(err)
			httpio.NotFound(w, err)        // что если httpio.ToWriter(w).WithLog(logger).NotFound(w)
			return
		}
		logger.Error().Err(err).Send()
		httpio.BadRequest(w, err)
		return
	}

	w.Header().Set("Content-Type", httpio.TypeTextPlain)
	w.Write([]byte(strconv.FormatInt(v, 10)))
}
