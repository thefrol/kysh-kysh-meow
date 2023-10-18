package retry

import (
	"errors"
	"fmt"
	"time"
)

type Options struct {
	delays     []time.Duration
	maxretries int
	conditions []func(error) bool
	callbacks  []RetryCallback
}

// This запускает несколько попыток запустить функцию,
// новые попытки будут совершаться, только ошибка
// обернута при помощи retry.Retriable(), или если выполняется
// одно из условий переданных при помощи If() или IfError()
func This(f func() error, opts ...Option) error {

	options := Options{}

	// укажем стандартные настройки
	Attempts(1)(&options)
	DelaySeconds(1)(&options)

	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return fmt.Errorf("ошибка опции при запуске восстанавливаемой функции: %w", err)
		}
	}

	var err error
	for i := 0; i <= options.maxretries; i++ {
		// Мы пускаем массив задержек по кругу, кроме первой попытки
		if i > 0 {
			currentI := int((i - 1) % len(options.delays))
			time.Sleep(options.delays[currentI])

			// вызываем коллбеки перед вызовом функции, то есть после того
			// как счетчик i поменял значение, а раньше нужно было бы проверить,
			// что следущая итерация случится
			for _, c := range options.callbacks {
				c(i+1, err) // i+1 потому что тут возвращаем человеческое значение, где попытки начинаются с 1
			}
			continue
		}

		err = f()
		if err == nil {
			return nil
		}
		if canretry(err, options) {
			// отсюда перенесен код, чтобы аглоритм легче читался
			continue
		}

		return err
	}
	return err
}

func canretry(err error, opts Options) bool {
	// todo
	//
	// Интересная тема, что в самой бы ошибке
	// можно было бы указать сколько раз функцию можно
	// перезапустить и с какими интервалами
	var re *RetriableError
	if errors.As(err, &re) {
		// запускаем коллбеки перед повторным запуском
		return true
	}

	for _, cond := range opts.conditions {
		if cond(err) {
			return true
		}
	}
	return false
}
