package handlers

import (
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
)

// ForJSON объединяет хендлеры, которые работают
// с джейсонами
type ForJSON struct {
	Registry manager.Registry
}
