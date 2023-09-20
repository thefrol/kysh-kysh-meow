package stats

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// DropPoll сбрасывает значение счетчика опросов памяти в указанном хранилище
// dropCounter сбраcывает счетчик
func DropPollCount(store storage.Storager) {
	dropCounter(store, metricPollCount)
}
func incrementPollCount(store storage.Storager) {
	incrementCounter(store, metricPollCount)
}

// incrementCounter обновляет PollCount счетчик в хранилище, добавляет ему единицу
func incrementCounter(store storage.Storager, name string) {
	count, _ := store.Counter(metricPollCount)
	store.SetCounter(metricPollCount, count+metrica.Counter(1))
}

// dropCounter сбраcывает счетчик
func dropCounter(store storage.Storager, name string) {
	_, found := store.Counter(metricPollCount)
	if !found {
		//если такого счетчика нет, то и не сбрасываем ничего
		return

	}
	store.SetCounter(name, metrica.Counter(0))
}
