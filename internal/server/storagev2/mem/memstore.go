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
	_ scan.Labler               = (*MemStore)(nil)
)

// MemStore хранит метрики в памяти, и сбрасывает
// их на диск
//
// Если указан FileStorage, то при каждом изменении счетчика
// счетчик будет записан в файл.
//
// Если нужна интервальная запись, то стоит воспользоваться
// IntervalicSaver
type MemStore struct {
	cmt      sync.RWMutex // мьютекс для счетчиков
	Counters IntMap

	gmt    sync.RWMutex // мьютекс для счетчиков
	Gauges FloatMap

	// Если указан FilePath, то сюда будут
	// синхронно писаться данные. При этом
	// Чтение просиходит не отсюда
	FilePath string

	Log zerolog.Logger
}

// All implements scan.CounterLister.
func (s *MemStore) Labels(context.Context) (map[string][]string, error) {
	// проверять, что мапа инициализирована, не нужно
	// но надо проверить, что сам MemStore не нулевой
	if s == nil {
		return nil, fmt.Errorf("MemStore: %w", ErrorNilStore)
	}

	const typeCount = 2 // это сколько типов метрик
	m := make(map[string][]string, typeCount)

	// займемся счетчиками

	cl := make([]string, 0, len(s.Counters))

	s.cmt.RLock()
	for c := range s.Counters {
		cl = append(cl, c)
	}
	s.cmt.RUnlock()

	m["counters"] = cl

	// теперь займемся гаужами

	gl := make([]string, 0, len(s.Gauges))

	s.gmt.RLock()
	for g := range s.Gauges {
		gl = append(gl, g)
	}
	s.gmt.RUnlock()

	m["gauges"] = gl

	//

	return m, nil
}
