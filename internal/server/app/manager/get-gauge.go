package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func (r Registry) Gauge(ctx context.Context, id string) (float64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", app.ErrorBadID, err)
	}

	return r.Gauges.Get(ctx, id)
}
