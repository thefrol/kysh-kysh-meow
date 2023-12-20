package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func (r Registry) UpdateGauge(ctx context.Context, id string, value float64) (float64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", app.ErrorBadID, err)
	}

	return r.Gauges.Update(ctx, id, value) // todo обернуть ошибку??
}
