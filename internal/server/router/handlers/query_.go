package handlers

import (
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
)

// ForQuery группирует в себе хендлеры для обработки базовых запросов,
// которые мы делали на первых инкрементах. Тут параметры берутся из
// URL запроса.
type ForQuery struct {
	Registry manager.Registry
}
