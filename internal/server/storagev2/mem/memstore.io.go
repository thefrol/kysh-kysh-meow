package mem

import (
	"fmt"
	"os"

	"github.com/mailru/easyjson"
)

type FileData struct {
	Counters IntMap
	Gauges   FloatMap
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
	d := FileData{
		Counters: s.Counters,
		Gauges:   s.Gauges,
	}

	s.cmt.RLock()
	s.gmt.RLock()
	buf, err := easyjson.Marshal(d)
	s.cmt.RUnlock()
	s.gmt.RUnlock()
	if err != nil {
		s.Log.Error().
			Err(err).
			Msg("Ошибка маршилинга")

		return fmt.Errorf("MemStore: %w", err)
	}

	// записываем в файл
	err = os.WriteFile(s.FilePath, buf, 0644)
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

	var d FileData
	if err = easyjson.Unmarshal(buf, &d); err != nil {
		s.Log.Error().
			Err(err).
			Msg("Ошибка демаршалинга данных в хранилище")

		return fmt.Errorf("MemStore: %w", err)
	}

	s.Log.Info().
		Err(err).
		Str("file", s.FilePath).
		Msg("хранилище прочитано")

	// заменяем каунтеры
	if s.Counters != nil {
		s.Log.Info().Msg("мапа с канутерами не пустая, и будет заменена")
	}
	s.Counters = d.Counters

	if s.Gauges != nil {
		s.Log.Info().Msg("мапа с гаужами не пустая, и будет заменена")
	}
	s.Gauges = d.Gauges

	return nil
}
