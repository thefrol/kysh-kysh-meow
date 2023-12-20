package app

import "errors"

var (
	ErrorBadID           = errors.New("неправильный id")
	ErrorMetricNotFound  = errors.New("метрика не найдена")
	ErrorUnknownMetric   = errors.New("неизвестный тип метрики")
	ErrorValidationError = errors.New("ошибка валидации")
)
