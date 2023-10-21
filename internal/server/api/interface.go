package api

import (
	"context"
	"errors"

	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
)

// Раз уж интерфейс тут, то, наверное и ошибки тоже должны быть описаны в этом же пакете
var (
	ErrorParseError        = errors.New("невозможно распарсить значение метрики в int или float")
	ErrorUnknownMetricType = errors.New("неизвестная метрика, доступны значения counter и gauge")
	ErrorUpdateCheckFailed = errors.New("обновление не удалось")
	ErrorNotFoundMetric    = errors.New("метрика с указанным именем не найдена")
	ErrorDeltaEmpty        = errors.New("поле Delta не может быть пустым, для когда id=counter")
	ErrorValueEmpty        = errors.New("поле Value не может быть пустым, для когда id=gauge")

	ErrorNoDatabase              = errors.New("база данных в текущей конфигурации не исползуется")
	ErrorNoConntectionToDatabase = errors.New("нет связи с базой данных")
)

type datastruct = metrica.Metrica

// Operation представляет собой операцию над хранилищем. Мы передаем такие операции в
// хендлеры
type Operation func(context.Context, ...datastruct) (out []datastruct, err error)

// Operator это новый интерфейс хранилища, его название само по себе дает намек. Оператор -
// в данном смысле это коллекция операций. Например, Get Update - мы хотим, чтобы они имели тип
// Operation, чтобы в дальнейшем можно было использовать с функциями HandleJSONRequest и ему
// подобным
type Operator interface {
	Get(ctx context.Context, req ...datastruct) (resp []datastruct, err error)
	Update(ctx context.Context, req ...datastruct) (resp []datastruct, err error)

	Check(ctx context.Context) error

	List(ctx context.Context) (counterNames []string, gaugeNames []string, err error)
}

// mentor
//
// Где должны храниться ошибки?

// TODO
//
// Внезапно пришла забавная идея.
//
// Если Operation это тип, то мы могли бы сделать к нему методов, которые и были бы этими хендлерами
// func (op Operation) HandleJSON
//
// И тогда бы наши вызовы выглядели вот так
//
// router.Post("/value", get.HadleWithJSON)
// router.Post("/update", update.HandleWithJSON)
//
// пока выглядит не очень идиоматично, канеш)
