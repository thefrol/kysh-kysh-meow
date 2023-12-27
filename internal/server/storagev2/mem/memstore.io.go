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
// если path не указан будет использован
// s.FilePath
func (s *MemStore) Dump(path string) error {
	if s == nil {
		return fmt.Errorf("запись в файл: %w", ErrorNilStore)
	}

	// если файл пустой, то ничего не пишем
	if path == "" {
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
	err = os.WriteFile(path, buf, 0644)
	if err != nil {
		s.Log.Error().
			Err(err).
			Str("file", path).
			Msg("Ошибка записи в файл")

		return fmt.Errorf("MemStore: %w", err)
	}

	s.Log.Debug().
		Str("file", path).
		Msg("Хранилище записано в файл")

	return nil
}

// Restore перезаписываем хранилище из файла path
// при этом не использует filePath из самой структуры
func (s *MemStore) RestoreFrom(path string) error {
	if s == nil {
		return fmt.Errorf("чтение из файла : %w", ErrorNilStore)
	}

	if path == "" {
		// не читаем
		return fmt.Errorf("файл для чтения не может быть пустым")
	}

	buf, err := os.ReadFile(path)
	if err != nil {
		s.Log.Error().
			Err(err).
			Str("file", path).
			Msg("ошибка чтения файла")
		return fmt.Errorf("MemStore: %w", err)
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
		Str("file", path).
		Msg("хранилище прочитано")

	// читаем из прочитанной мапы d и записываем в свою
	s.cmt.Lock()
	// возможно мапа пустая, тогда
	// создаем новую мапу, но она будет
	// в куче, а чтобы этого не случилсоь,
	// её нужно создать заранее
	if s.Counters == nil {
		s.Counters = make(IntMap)
	}

	for i, c := range d.Counters {
		if _, ok := s.Counters[i]; ok {
			s.Log.Info().
				Str("id", i).
				Str("file", path).
				Msg("счетчик будет переписан из файла")
		}
		s.Counters[i] = c
	}
	s.cmt.Unlock()

	// теперь то же самое с гаужами сделаем
	s.gmt.Lock()
	// возможно мапа пустая, тогда
	// создаем новую мапу, но она будет
	// в куче, а чтобы этого не случилсоь,
	// её нужно создать заранее
	if s.Gauges == nil {
		s.Gauges = make(FloatMap)
	}
	for i, g := range d.Gauges {
		if _, ok := s.Gauges[i]; ok {
			s.Log.Info().
				Str("id", i).
				Str("file", path).
				Msg("гауж будет переписан из файла")
		}
		s.Gauges[i] = g
	}
	s.gmt.Unlock()

	return nil
}
