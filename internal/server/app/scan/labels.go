// Этот пакет занимается тем, что выводит все сохраненные метрики.
// По сути это один из юз-кейсов, которые я не хочу перемешивать
// с остальными.
//
// Тут мы выдаем данные всех метрик, которые у нас есть в базе.
//
// Все метрики нам нужны, чтобы вывести HTML страницу, и может быть
// в будущем чтобы передать все эти метрики в Прометеус
package scan

import (
	"context"
	"fmt"
)

type Labler interface {
	Labels(context.Context) (map[string][]string, error)
}

// Labels используется чтобы вернуть имена всех
// метрик, которые мы знаем
type Labels struct {
	Labels Labler
}

// Get возващет все имена метрик, собранные по типам метрик
func (l Labels) Get(ctx context.Context) (map[string][]string, error) {
	ls, err := l.Labels.Labels(ctx)
	if err != nil {
		return nil, fmt.Errorf("labels: %w", err)
	}

	return ls, nil
}
