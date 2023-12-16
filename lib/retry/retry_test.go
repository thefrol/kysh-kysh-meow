package retry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thefrol/kysh-kysh-meow/lib/retry"
)

func Test_RetriableErrorWorks(t *testing.T) {
	t.Run("Обернутость в ретриабл работает", func(t *testing.T) {
		callee := func() error {
			return retry.Retriable(errors.New("test error"))
		}

		start := time.Now()
		retry.This(callee,
			retry.Attempts(1),
			retry.DelaySeconds(1, 1, 2))

		dur := time.Since(start)
		assert.True(t, dur > time.Millisecond*1000, "Должно хотя бы секунду длиться")
	})

	t.Run("Если ошибка не ретриабл, то не повторяем", func(t *testing.T) {
		callee := func() error {
			return errors.New("not retriable")
		}

		start := time.Now()
		retry.This(callee,
			retry.Attempts(2),
			retry.DelaySeconds(1, 1, 2))

		dur := time.Since(start)
		assert.True(t, dur < time.Millisecond*10, "Не должно быть задержек")
	})

}

func Test_Callbacks(t *testing.T) {
	t.Run("четыре запуска, три ретрая, три коллбека", func(t *testing.T) {

		callee := func() error {
			return retry.Retriable(errors.New("test error"))
		}

		counter := 0
		increment := func(int, error) {
			counter++
		}

		retry.This(callee,
			retry.Attempts(3),
			retry.DelaySeconds(1, 1, 1),
			retry.OnRetry(increment))

		assert.Equal(t, 3, counter, "Коллбеки должны были запуститься три раза")
	})
}

func Test_CountOfRuns(t *testing.T) {
	t.Run("четыре запуска, четыре прохода по сновной функции", func(t *testing.T) {

		counter := 0
		callee := func() error {
			counter++
			return retry.Retriable(errors.New("test error"))
		}

		retry.This(callee,
			retry.Attempts(3),
			retry.DelaySeconds(1, 1, 1))

		assert.Equal(t, 4, counter, "Коллбеки должны были запуститься три раза")
	})
}
