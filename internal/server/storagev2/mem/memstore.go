package mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
)

var (
	_ manager.CounterRepository = (*MemStore)(nil)
	_ manager.GaugeRepository   = (*MemStore)(nil)
	_ scan.CounterLister        = (*MemStore)(nil)
)

type MemStore struct {
	cmt      sync.RWMutex // мьютекс для счетчиков
	Counters IntMap

	gmt    sync.RWMutex // мьютекс для счетчиков
	Gauges FloatMap

	Log      zerolog.Logger
	FilePath string
}

// All implements scan.CounterLister.
func (s *MemStore) All(context.Context) (map[string]int64, error) {
	// проверять, что мапа инициализирована, не нужно
	// но надо проверить, что сам MemStore не нулевой
	if s == nil {
		return nil, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	m := make(map[string]int64, len(s.Counters))

	// заблокируем
	s.cmt.Lock()
	defer s.cmt.Unlock()

	// скопируем мапу в новую мапу
	for k, v := range s.Counters {
		m[k] = v
	}

	// вернем копию
	return m, nil
}
