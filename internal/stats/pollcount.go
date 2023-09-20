package stats

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// incrementCounter обновляет PollCount счетчик в хранилище, добавляет ему единицу
func IncrementCounter(store storage.Storager, name string) {
	count, _ := store.Counter(metricPollCount)
	store.SetCounter(metricPollCount, count+metrica.Counter(1))
}

// incrementCounter сбравыем счетчик
func DropCounter(store storage.Storager, name string) {
	_, found := store.Counter(metricPollCount)
	if !found {
		//если такого счетчика нет, то и не сбрасываем ничего
		return

	}
	store.SetCounter(name, metrica.Counter(0))
}
