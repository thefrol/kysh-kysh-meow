package mem

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// тестируем работу с неинициализированными стораджами
func Test_NoPanic(t *testing.T) {

	ctx := context.Background()

	t.Run("можем прочитать из nil мемстора, без паники", func(t *testing.T) {
		// инициализируем именно ссылку нулевую
		var ms *MemStore

		ms.Counter(ctx, "yo")
	})

	t.Run("можем записать в nil MemStore, будет ошибка но без паники", func(t *testing.T) {
		// инициализируем именно ссылку нулевую
		var ms *MemStore

		_, err := ms.CounterIncrement(ctx, "yo", 2)
		assert.ErrorIs(t, err, ErrorNilStore)
	})

	t.Run("можем записать без инициализации мап", func(t *testing.T) {
		var ms MemStore

		const (
			value = 2
			id    = "yo"
		)

		// мапу не инициализировали, но запишем в нее
		_, err := ms.CounterIncrement(ctx, id, value)
		assert.NoError(t, err, ErrorNilStore)

		// и прочитаем!
		v, err := ms.Counter(ctx, id)
		assert.NoError(t, err, ErrorNilStore)
		assert.EqualValues(t, value, v)
	})
}

func Test_FileSave(t *testing.T) {
	// тут нужен сюите, чтобы удалять файлы
	t.Run("файл создается", func(t *testing.T) {
		var (
			file = "./file.test"
		)

		err := os.Remove(file)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Не могу удалить тестовый файл %s", file)
		}

		// удостоверимся, что файла реально нет
		_, err = os.Stat(file)
		assert.ErrorIs(t, err, os.ErrNotExist)

		var ms = MemStore{
			FilePath: file,
		}

		// записываем хранилище в файл
		err = ms.Dump()
		assert.NoError(t, err)

		_, err = os.Stat(file)
		assert.NoError(t, err)

		err = os.Remove(file)
		if err != nil && err != os.ErrNotExist {
			t.Errorf("Не могу удалить тестовый файл %s", file)
		}
	})

	t.Run("чтение из несуществующего файла дает ошибку", func(t *testing.T) {
		var (
			file = "./file.test"
		)

		err := os.Remove(file)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Не могу удалить тестовый файл %s", file)
		}

		// удостоверимся, что файла реально нет
		_, err = os.Stat(file)
		assert.ErrorIs(t, err, os.ErrNotExist)

		var ms = MemStore{
			FilePath: file,
		}

		// записываем хранилище в файл
		err = ms.Restore()
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("запишем и прочитаем", func(t *testing.T) {
		t.Error("not implemented")
	})
}

func Test_Counters(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name   string
		values []int64
		sum    int
		id     string
	}{
		{
			name:   "только положительные",
			values: []int64{1, 3, 4},
			sum:    8,
			id:     "yo",
		},
		{
			name:   "и отрицательные",
			values: []int64{1, -3, 4},
			sum:    2,
			id:     "yo",
		},
		{
			name:   "и нули",
			values: []int64{1, 0, -1},
			sum:    0,
			id:     "yo",
		},
		{
			name:   "только нули",
			values: []int64{0, 0, 0},
			sum:    0,
			id:     "yo",
		},
	}

	for _, tt := range tests {
		t.Run("Инкрементируем только положительные", func(t *testing.T) {
			var ms MemStore

			// мапу не инициализировали, но запишем в нее
			for _, v := range tt.values {
				_, err := ms.CounterIncrement(ctx, tt.id, v)
				assert.NoError(t, err, ErrorNilStore)

			}

			// и прочитаем!
			v, err := ms.Counter(ctx, tt.id)
			assert.NoError(t, err, ErrorNilStore)
			assert.EqualValues(t, tt.sum, v)
		})

	}
}
