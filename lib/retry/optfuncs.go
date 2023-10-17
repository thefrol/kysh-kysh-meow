package retry

import (
	"errors"
	"time"
)

type Option func(opt *Options) error

// DelaySeconds позволяет указать задержки между попытками,
// если повторов будет больше, чем задержек, то задержки
// будут повторяться по кругу
func DelaySeconds(seconds ...uint) Option {
	return func(opt *Options) error {
		var delays []time.Duration
		for _, interval := range seconds {
			delays = append(delays, time.Second*time.Duration(interval))
		}
		opt.delays = delays
		return nil
	}
}

func OnRetry(funs ...func(int)) Option {
	return func(opt *Options) error {
		opt.callbacks = append(opt.callbacks, funs...)
		return nil
	}
}

// Attempts позволяет установить количество повторных попыток,
// и это значит что указав Attempts(2), мы запустим функцию ровно три раза
func Attempts(count uint) Option {
	return func(opt *Options) error {
		opt.maxretries = int(count) + 1
		return nil
	}
}

// If позволяет указать функцию, которая проверит можно ли с
// ошибкой err делать ещё одну попытку.
//
// В одном запуске может быть несколько таких условий, функция сработает
// если выполнится одно из них
func If(f func(error) bool) Option {
	return func(opt *Options) error {
		opt.conditions = append(opt.conditions, f)
		return nil
	}
}

// IfError позволяет повторить запуск, если повторяемая функция вернула
// определенную ошибку, например io.EOF
//
// В одном запуске может быть несколько таких условий, функция сработает
// если выполнится одно из них
func IfError(target error) Option {
	return func(opt *Options) error {
		opt.conditions = append(opt.conditions, func(err error) bool {
			return errors.Is(err, target)
		})
		return nil
	}
}
