package domain

import "errors"

var (
	ErrorBadID          = errors.New("неправильный id")
	ErrorMetricNotFound = errors.New("метрика не найдена")
)
