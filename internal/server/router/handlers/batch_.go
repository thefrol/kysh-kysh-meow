package handlers

import "github.com/thefrol/kysh-kysh-meow/internal/server/app/metricas"

// ForBatch объединяет хендлеры, которые работают
// с пачками джейсонов
//
// Вообще почему я выделил это в отдельный юзкейс:
// потому что тут нужны будут транзакции. Транзакции это
// вообще другие интерфейсы и вообще другой способ работать
// с обновлениями
type ForBatch struct {
	Manager metricas.Manager
}
