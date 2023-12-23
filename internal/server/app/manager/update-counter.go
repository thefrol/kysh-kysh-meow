package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

func (r Registry) IncrementCounter(ctx context.Context, id string, delta int64) (int64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", app.ErrorBadID, err)
	}

	return r.Counters.CounterIncrement(ctx, id, delta) // todo обернуть ошибку??
}
