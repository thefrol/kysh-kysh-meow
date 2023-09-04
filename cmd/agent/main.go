package main

import (
	"fmt"

	"github.com/thefrol/kysh-kysh-meow/internal/storage"
)

var store storage.Storager

func init() {
	store = storage.New()
}

func main() {

	saveMemStats(store, nil)
	fmt.Println(store.Gauge("Alloc"))
	fmt.Println(store.Gauge("HeapSys"))
}
