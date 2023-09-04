package main

import (
	"fmt"
	"runtime"
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
	fields, err := getStructFields(m)
	if err != nil {
		return nil, nil, fmt.Errorf("can't retrieve fields from MemStats")
	}
	lost = Difference[string](metricsMem, fields)
	exclude = Difference[string](fields, metricsMem)
	return lost, exclude, nil
}
