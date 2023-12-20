package mem

import (
	"fmt"
	"os"

	"github.com/mailru/easyjson"
)

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
