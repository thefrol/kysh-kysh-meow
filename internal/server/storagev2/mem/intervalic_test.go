package mem

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartAndStop(t *testing.T) {
	var s = IntervalicSaver{
		Store:    &MemStore{},
		File:     "dummy.test",
		Interval: 100 * time.Second,
	}
	err := s.Run()
	require.NoError(t, err, "не удалось запустить хранилише")

	err = s.Stop()
	assert.NoError(t, err, "не удалось остановить хранилище")

	// файл должен создасться с конце, хоть там ничего может и не быть
	// потому что в конце мы все равно записываем последние изменения
	_, err = os.Stat(s.File)
	assert.NoError(t, err, "файл должен создасться путь мы в него и ничего не писали")

	err = os.Remove(s.File)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Не могу удалить тестовый файл %s", s.File)
	}
}
func TestEmptyFile(t *testing.T) {
	var s = IntervalicSaver{
		Store:    &MemStore{},
		File:     "",
		Interval: 100 * time.Second,
	}

	err := s.Run()
	require.ErrorIs(t, err, ErrorBadConfig)
	assert.False(t, s.started)
}

func TestNilStore(t *testing.T) {
	var s = IntervalicSaver{
		Store:    nil,
		File:     "dummy",
		Interval: 100 * time.Second,
	}

	err := s.Run()
	require.ErrorIs(t, err, ErrorBadConfig)
	assert.False(t, s.started)
}

func TestStopNotStarted(t *testing.T) {
	var s IntervalicSaver

	err := s.Stop()
	assert.Error(t, err)
	assert.False(t, s.started)
}
