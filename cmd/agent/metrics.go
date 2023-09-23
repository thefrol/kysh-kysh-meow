package main

import (
	"fmt"
	"runtime"

	"github.com/thefrol/kysh-kysh-meow/internal/slices"
	"github.com/thefrol/kysh-kysh-meow/internal/structs"
)

var metricsMem = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// parseMetrics возвращает метрики, не найденные в выдаче runtime(lost) и список несохраняемых метрик(excluded).
// В случае ошибки вернет ошибку третим параметром
func parseMetrics() (lost []string, exclude []string, err error) {
	m := runtime.MemStats{} // по хорошему тут бы получать какой-то пустой элемент, и чтобы getFields работала только с типом!
	fields, err := structs.FieldNames(m)
	if err != nil {
		return nil, nil, fmt.Errorf("can't retrieve fields from MemStats")
	}
	lost = slices.Difference[string](metricsMem, fields)
	exclude = slices.Difference[string](fields, metricsMem)
	return lost, exclude, nil
}
