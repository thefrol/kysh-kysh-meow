package main

import (
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

// incrementCounter обновляет PollCount счетчик в хранилище, добавляет ему единицу
func incrementCounter(store storage.Storager, name string) error {
	count, _ := store.Counter(metricPollCount)
	store.SetCounter(metricPollCount, count+metrica.Counter(1))
	return nil
}

// incrementCounter сбравыем счетчик
func dropCounter(store storage.Storager, name string) error {
	store.SetCounter(name, metrica.Counter(0))
	return nil
}
