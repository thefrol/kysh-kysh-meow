package mem

import (
	"context"
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
)

// Counter implements manager.CounterRepository.
func (s *MemStore) Counter(ctx context.Context, id string) (int64, error) {
	// проверять, что мапа инициализирована, не нужно
	// но надо проверить, что сам MemStore не нулевой
	if s == nil {
		return 0, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	// если в мапе не нашли значения,
	// то возврашаем, что метрика не найдена
	s.cmt.RLock()
	v, ok := s.Counters[id]
	s.cmt.RUnlock()
	if !ok {
		return 0, app.ErrorMetricNotFound
	}

	return v, nil
}

// CounterIncrement implements manager.CounterRepository.
func (s *MemStore) CounterIncrement(ctx context.Context, id string, delta int64) (int64, error) {
	// но надо проверить, что сам MemStore не нулевой
	// и мапа
	if s == nil {
		return 0, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	// если мапа не создана - создаем,
	// тут можно было бы подумать над
	// емкостью этой мапы изначальной
	if s.Counters == nil {
		s.Counters = make(IntMap)
	}

	// нам все равно есть
	// или нет такой метрики
	// получается как будто
	// она равна нулю
	s.cmt.Lock()
	s.Counters[id] += delta
	s.cmt.Unlock()

	// запишем обновления в файл
	s.Dump("")

	return s.Counters[id], nil
}
