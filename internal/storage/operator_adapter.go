package storage

import (
	"context"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
)

// OperatorAdapter оборачивает классы на старом апи хранилища oldAPI,
// под новое апи api.Storager, делая так, контекст, конечно используется тупо вхолостую
type OperatorAdapter struct {
	legacyStore legacyStorager
}

// AsOperator оборачивает хранилища старого интерфейса, позволяя их подключать к хендлерам, работающим на новом интерфейсе api.Storager
func AsOperator(s legacyStorager) *OperatorAdapter {
	return &OperatorAdapter{legacyStore: s}
}

// Counter implements api.Storager.
func (a OperatorAdapter) Counter(ctx context.Context, name string) (value int64, err error) {
	c, found := a.legacyStore.Counter(name)
	if !found {
		return 0, api.ErrorNotFoundMetric
	}
	return int64(c), nil
}

// Gauge implements api.Storager.
func (a OperatorAdapter) Gauge(ctx context.Context, name string) (value float64, err error) {
	g, found := a.legacyStore.Gauge(name)
	if !found {
		return 0, api.ErrorNotFoundMetric
	}
	return float64(g), nil
}

// IncrementCounter implements api.Storager.
func (a *OperatorAdapter) IncrementCounter(ctx context.Context, name string, delta int64) (value int64, err error) {
	was, _ := a.legacyStore.Counter(name)
	newVal := was + metrica.Counter(delta)
	a.legacyStore.SetCounter(name, newVal)
	return int64(newVal), nil
}

// UpdateGauge implements api.Storager.
func (a *OperatorAdapter) UpdateGauge(ctx context.Context, name string, v float64) (value float64, err error) {
	a.legacyStore.SetGauge(name, metrica.Gauge(v))
	return v, nil
}

// List implements api.Storager.
func (a *OperatorAdapter) List(ctx context.Context) (counterNames []string, gaugeNames []string, err error) {
	return a.legacyStore.ListCounters(), a.legacyStore.ListGauges(), nil
}

func (a *OperatorAdapter) getOne(ctx context.Context, r metrica.Metrica) (resp metrica.Metrica, err error) {
	if r.ID == "" {
		return resp, api.ErrorUpdateCheckFailed // todo правильная ошибка?
	}
	switch r.MType {
	case "counter":
		c, err := a.Counter(ctx, r.ID)
		return metrica.Metrica{MType: r.MType, ID: r.ID, Delta: &c}, err
	case "gauge":
		g, err := a.Gauge(ctx, r.ID)
		return metrica.Metrica{MType: r.MType, ID: r.ID, Value: &g}, err
	default:
		return resp, api.ErrorUnknownMetricType
	}
}

func (a *OperatorAdapter) updateOne(ctx context.Context, in metrica.Metrica) (resp metrica.Metrica, err error) {

	switch in.MType {
	case "counter":
		if in.Delta == nil {
			return empty, api.ErrorDeltaEmpty
		}
		c, err := a.IncrementCounter(ctx, in.ID, *in.Delta)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Delta: &c}, err // это получается отправится в хип
	case "gauge":
		if in.Value == nil {
			return empty, api.ErrorValueEmpty
		}
		g, err := a.UpdateGauge(ctx, in.ID, *in.Value)
		return metrica.Metrica{MType: in.MType, ID: in.ID, Value: &g}, err
	default:
		return empty, api.ErrorUnknownMetricType
	}
}

// TODO мне конечно очень не нравится что мы используем траспортный класс для работы внутри хранилища,
// но пока в учебных целях так. Будем надеятся все сойдется. Просто если делать ещё один класс, там надо будет
// переобертывание писать итд.
//
// Решение на пока - псевдоним типа DataStruct

type datastruct = metrica.Metrica

type Operator func(context.Context, datastruct) (datastruct, error)

// aggregate позволяет обрабатывать пачку переменных, имея только функцию которая работает с единственной переменной.
// Так, например, getOne принимает себе единственную metrica.Metrica, а с помощью функции aggregate мы можем сразу
// обрабаоть целую кучу Metrica при поможи getOne
//
// GetMany:
// many,err:= aggregate(getOne, ctx, batch...)
func aggregate(o Operator, ctx context.Context, rs ...datastruct) (resp []datastruct, err error) {
	for _, in := range rs {
		out, err := o(ctx, in)
		if err != nil {
			return resp, err
		}
		resp = append(resp, out)
	}
	return
}

// Get implemests Operator
func (a *OperatorAdapter) Get(ctx context.Context, req ...datastruct) (resp []datastruct, err error) {
	return aggregate(a.getOne, ctx, req...)
}

// Update implemests Operator
func (a *OperatorAdapter) Update(ctx context.Context, req ...datastruct) (resp []datastruct, err error) {
	return aggregate(a.updateOne, ctx, req...)
}

// Check implements Operator. Это функция, которая проверяет связь с базой данных. Возвращает ошибку если
// связь отсутствует
func (a *OperatorAdapter) Check(ctx context.Context) error {
	return api.ErrorNoDatabase

	// TODO
	//
	// Мне кажется это немного неправильно иметь такой метод.
	// По сути все классы хранилища теперь должны знать о базе данных,
	// Но и прокидывать базу данных как-то в обход не хочется пока.
	//
	// Кажется, что это меньшее из зол.
}

var _ api.Operator = (*OperatorAdapter)(nil)

var empty = metrica.Metrica{}
