package mem

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

// Gauge implements manager.GaugeRepository.
func (s *MemStore) Gauge(ctx context.Context, id string) (float64, error) {
	// проверять, что мапа инициализирована, не нужно
	// но надо проверить, что сам MemStore не нулевой
	if s == nil {
		return 0, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	// если в мапе не нашли значения,
	// то возврашаем, что метрика не найдена
	s.gmt.RLock()
	v, ok := s.Gauges[id]
	s.gmt.RUnlock()
	if !ok {
		return 0, app.ErrorMetricNotFound
	}

	return v, nil
}

// GaugeIncrement implements manager.GaugeRepository.
func (s *MemStore) GaugeUpdate(ctx context.Context, id string, value float64) (float64, error) {
	// но надо проверить, что сам MemStore не нулевой
	// и мапа
	if s == nil {
		return 0, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	// если мапа не создана - создаем,
	// тут можно было бы подумать над
	// емкостью этой мапы изначальной
	if s.Gauges == nil {
		s.Gauges = make(FloatMap)
	}

	// нам все равно есть
	// или нет такой метрики
	// получается как будто
	// она равна нулю
	s.gmt.Lock()
	s.Gauges[id] = value
	s.gmt.Unlock()

	// запишем обновления в файл
	s.Dump()

	return s.Gauges[id], nil
}
