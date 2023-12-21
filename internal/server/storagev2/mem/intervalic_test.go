package mem

import (
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
