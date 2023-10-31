package collector

import (
	"context"
)

const MaxBatch = 400

func pool(ctx context.Context, count int, worker func()) {

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func() {
		loop:
			for {
				select {
				case <-ctx.Done():
					break loop
				default:
					worker()
				}
			}
			wg.Done()
		}()

	}
}
