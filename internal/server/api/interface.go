package api

import "context"

// Storager это интерфейс к хранилищу, которое использует именно этот API. Таким образом мы делаем хранилище зависимым от
// API,  а не наоборот
type Storager interface {
	Counter(ctx context.Context, name string) (value int64, found bool, err error)
	UpdateCounter(ctx context.Context, name string, delta int64) (newValue int64, err error)
	ListCounters(ctx context.Context) ([]string, error)
	Gauge(ctx context.Context, name string) (value float64, found bool, err error)
	UpdateGauge(ctx context.Context, name string, value float64) (newValue float64, err error)
	ListGauges(ctx context.Context) ([]string, error)
}

// todo
//
// КАжется в дальнейшем это надо вынести ещё выше в каталогах
