package mem

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		err = ms.Dump(ms.FilePath) // todo нужен какой-то рефактиринг этой функции, она должна быть отдельно от стоража
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
		err = ms.RestoreFrom(file)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("запишем и прочитаем", func(t *testing.T) {
		var (
			file = "./file.test"

			gaugeID    = "g"
			gaugeValue = 3

			counterID    = "c"
			counterValue = 6
		)

		ctx := context.Background()

		err := os.Remove(file)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Не могу удалить тестовый файл %s", file)
		}

		// удостоверимся, что файла реально нет
		_, err = os.Stat(file)
		assert.ErrorIs(t, err, os.ErrNotExist)

		// создадим первый стораж,
		// который пишет в файл
		// и запишем в файл

		var ms = MemStore{
			FilePath: file,
		}

		_, err = ms.CounterIncrement(ctx, counterID, int64(counterValue))
		assert.NoError(t, err, "не удалось записать счетчик")

		_, err = ms.GaugeUpdate(ctx, gaugeID, float64(gaugeValue))
		assert.NoError(t, err, "не удалось записать гауж")

		defer func() {
			// почистим за собой
			err = os.Remove(file)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				t.Errorf("Не могу удалить тестовый файл %s", file)
			}
		}()

		// теперь создаем второй стораж,
		// и прочитаем из файла

		var ls = MemStore{
			FilePath: file,
		}

		err = ls.RestoreFrom(file)
		require.NoError(t, err, "ошибка загрузки из файла")

		// теперь проверим что мы получили
		c, err := ls.Counter(ctx, counterID)
		assert.NoError(t, err, "не найден счетчик после загрузки")
		assert.EqualValues(t, counterValue, c)

		g, err := ls.Gauge(ctx, gaugeID)
		assert.NoError(t, err, "не найден гауж после загрузки")
		assert.EqualValues(t, gaugeValue, g)

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
		t.Run(tt.name, func(t *testing.T) {
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

// протестируем гаужи
func Test_Gauges(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name   string
		values []float64
		got    float64
		id     string
	}{
		{
			name:   "набор 1",
			values: []float64{1, 3, 4},
			got:    4,
			id:     "yo",
		},

		{
			name:   "набор 2",
			values: []float64{1, 3, -4.0000000001},
			got:    float64(-4.0000000001),
			id:     "yo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ms MemStore

			// мапу не инициализировали, но запишем в нее
			for _, v := range tt.values {
				_, err := ms.GaugeUpdate(ctx, tt.id, v)
				assert.NoError(t, err, ErrorNilStore)

			}

			// и прочитаем!
			v, err := ms.Gauge(ctx, tt.id)
			assert.NoError(t, err, ErrorNilStore)
			assert.EqualValues(t, tt.got, v)
		})

	}
}

// протестируем получение названий счетчиков
func Test_Labels(t *testing.T) {
	var m = MemStore{
		Counters: make(IntMap),
		Gauges:   make(FloatMap),
	}

	var (
		// метрики, которые мы установим
		gaugeID   = "g"
		counterID = "c"

		// количество гаужей и счетчиков по итогу тестов
		gaugeCount   = 1
		counterCount = 1
	)

	m.Counters[counterID] = 1
	m.Gauges[gaugeID] = 1

	// эту функцию мы и тестируем
	l, err := m.Labels(context.Background())
	require.NoError(t, err)

	// в результатах должны быть наши метрики
	assert.Contains(t, l["gauges"], gaugeID, "В гаужаж не найдена проставленная метрика")
	assert.Contains(t, l["counters"], counterID, "в счетчиках не найдена проставленна метрика")

	// и ничего лишнего
	assert.Equal(t, gaugeCount, len(l["gauges"]), "полученное количество гаужей должено быть другим")
	assert.Equal(t, counterCount, len(l["counters"]), "полученное количество счетчиков должено быть другим")
}
