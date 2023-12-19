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

type CounterLister interface {
	All(context.Context) (map[string]int64, error)
}

type GaugeLister interface {
	All(context.Context) (map[string]float64, error)
}

// Labels используется чтобы вернуть имена всех
// метрик, которые мы знаем
type Labels struct {
	Counters CounterLister
	Gauges   GaugeLister
}

// Get возващет все имена метрик, собранные по типам метрик
func (l Labels) Get(ctx context.Context) (map[string][]string, error) {
	res := make(map[string][]string)

	cs, err := l.Counters.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("in Labels.Get() Не удалось получить список cчетчиков %w", err)
	}

	arr := make([]string, 0, len(cs))
	for k := range cs {
		arr = append(arr, k)
	}
	res["counters"] = arr

	gs, err := l.Gauges.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("in Labels.Get() Не удалось получить список гаужей %w", err)
	}

	arr = make([]string, 0, len(gs))
	for k := range cs {
		arr = append(arr, k)
	}
	res["gauges"] = arr

	return res, nil
}
