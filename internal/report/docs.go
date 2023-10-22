// Пакет занимается сбором и отправкой метрик на сервер. А так же незаметно регулирует значение счетчика pollCount.
// Значение большинства метрик забираеются из пакета runtime.ReadMemStats() и ещё дополнительно пополняются двумя параметрами,
// один из которых счетчик PollCount - там хранится число раз, сколько мы опросили память. После отправки на сервер данных, это
// значение сбрасывается
//
// Основые функции:
//
// Fetch() - собирает основые параметры использования памяти и сохраняет
// в промеждуточное хранилише типа Stats, а так же увеличивает счетчик PollCount
//
// Send() - отправляет ранее собранные метрики на сервер, и обнуляет значение счетчика PollCount
// stats собирает информацию о использовании памяти, и сохраняет в хранилище
// так же имеет методы для управления счетчиком количества опросов PollCount
//
// Типичное использование:
//
// s := report.Fetch()
// err := report.Send(s.ToTransport(),"http://my.serv.er/update")
//
//	if err != nil{
//		...
//	}
//
// Имеет встроенное сжатие gzip, так что не надо об этом беспокоится
package report

// TODO
//
// Теперь мне не нравится, что тут больно дофига всего замешано)
// Тут и какие-то мидлвари, и опрос памяти и отправка
// Мидлварь, должна лежать наверное в агенте все же
