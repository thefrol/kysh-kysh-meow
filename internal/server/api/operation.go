package api

type Data struct{}

type Operation func(in []Data) (out []Data, err error)

// Стоит подумать о таком исполнении, это тогда будет центральным образом
// У нас есть какие-то опеации с хранилищем это Update, Get
//
// тогда объявление маршрута будет выглядеть вот так
//
// r.POST("/route", wrappers.JSON(operations.Get(store)))
//
// в operations будет лежать логика с нашими метриками
//
// func Get(store Storager) OperationFunc{}
//
// func JSON(op OperationFunc) http.HandlerFunc
//
//
// как-то так
