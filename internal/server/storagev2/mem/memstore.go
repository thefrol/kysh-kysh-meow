package mem

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/scan"
)

var (
	_ manager.CounterRepository = (*MemStore)(nil)
	_ scan.CounterLister        = (*MemStore)(nil)
)

type MemStore struct {
	cmt      sync.RWMutex // мьютекс для счетчиков
	Counters IntMap

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
	s.Dump()

	return s.Counters[id], nil
}

// Dump сохраняет хранилище в файл
func (s *MemStore) Dump() error {
	if s.FilePath == "" {
		// не пишем в файл
		return nil
	}

	// запишем джесон в файл
	// блокируем вот так, чтобы не мешать другим процессам,
	// пока мы уже в файл будем писать
	s.cmt.RLock()
	buf, err := easyjson.Marshal(s.Counters)
	s.cmt.RUnlock()
	if err != nil {
		s.Log.Error().
			Err(err).
			Msg("Ошибка маршилинга")

		return fmt.Errorf("MemStore: %w", err)
	}

	// записываем в файл
	err = os.WriteFile(s.FilePath, buf, os.ModeExclusive)
	if err != nil {
		s.Log.Error().
			Err(err).
			Str("file", s.FilePath).
			Msg("Ошибка записи в файл")

		return fmt.Errorf("MemStore: %w", err)
	}

	s.Log.Debug().
		Str("file", s.FilePath).
		Msg("Хранилище записано в файл")

	return nil
}

// Restore перезаписываем хранилище из файла
func (s *MemStore) Restore() error {
	if s.FilePath == "" {
		// не пишем в файл
		return nil
	}

	buf, err := os.ReadFile(s.FilePath)
	if err != nil {
		s.Log.Error().
			Err(err).
			Str("file", s.FilePath).
			Msg("ошибка чтения файла")
		return fmt.Errorf("MemStore: %w", err)
	}

	if err != nil {
		return fmt.Errorf("mem: %w", err)
	}

	s.cmt.Lock()
	err = easyjson.Unmarshal(buf, &s.Counters) // todo нужно проверить как это будет себя вести, когда тут мы перезаписываем, а там читаем. Что и куда он перезапишет)
	s.cmt.Unlock()
	if err != nil {
		s.Log.Error().
			Err(err).
			Msg("Ошибка демаршалинга данных в хранилище")

		return fmt.Errorf("MemStore: %w", err)
	}

	s.Log.Info().
		Err(err).
		Str("file", s.FilePath).
		Msg("хранилище прочитано")

	return nil
}
