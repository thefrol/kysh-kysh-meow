package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (r Registry) UpdateGauge(ctx context.Context, id string, value float64) (float64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", domain.ErrorBadID, err)
	}

	return r.Gauges.Update(ctx, id, value) // todo обернуть ошибку??
}
