package handlers

import "github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"

// ForJSON объединяет хендлеры, которые работают
// с джейсонами
type ForJSON struct {
	Manager metricas.Manager
}
