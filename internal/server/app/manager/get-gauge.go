package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (r Registry) Gauge(ctx context.Context, id string) (float64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", domain.ErrorBadID, err)
	}

	return r.Gauges.Get(ctx, id)
}
