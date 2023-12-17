package manager

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/domain"
)

func (r Registry) IncrementCounter(ctx context.Context, id string, delta int64) (int64, error) {
	if err := IsValid(id); err != nil {
		return 0, fmt.Errorf("%w: %v", domain.ErrorBadID, err)
	}

	return r.Counters.Increment(ctx, id, delta) // todo обернуть ошибку??
}
