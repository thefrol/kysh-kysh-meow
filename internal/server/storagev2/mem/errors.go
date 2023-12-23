package mem

import (
	"errors"
)

var (
	ErrorNotStarted = errors.New("не запущено")
	ErrorBadConfig  = errors.New("неправильно установлено")
)
