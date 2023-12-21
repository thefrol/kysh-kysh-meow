package mem

import (
	"errors"
)

var (
	ErrorNilRef     = errors.New("нулевая ссылка")
	ErrorNotStarted = errors.New("не запущено")
	ErrorBadConfig  = errors.New("неправильно установлено")
)
