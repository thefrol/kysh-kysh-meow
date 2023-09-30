// Пакет журнарирования Ололог
//
// Представляет собой синглтон-пакет. Можно пользоваться функциями напрямую:
// + ololog.Warn().Msg("тревожное сообщение")
// + ololog.Debug().Msg("дрюжелюбное сообщение как все хорошо")
//
// или через эксземпляр журнала
// ololog.Log().Warn("предупреждаю (^ж^)")
//
// Предоставляет отправку форматированных данных, в формате джейсон, например так:
//
//	ololog.Info().Str("band","Tame Impale").Int(rating,10).Send()
//
// на выходе отдаст
//
//	{"level":"info","band":"Tame Impala","rating":10,"time":1695268557,"message":"Крутейшая банда"}
//
// Так же любой пакет может создать некий поджурнал с нужными параметрами вот так
//
//	var localLog=ololog.Log().With().Str("foo", "bar").Logger()
//	localLog.Info().Msg("Info message") // {"level":"info", "location":"handlers","message":"Info message"}
package ololog

// Вообще, конечно хотелось бы, чтобы этот пакет имел такие функции
// ololog.UseZeroLog()
// ololog.UserZapLog()
//
// Для этого нужно лишь написать некую обертку для каждого такого журнала.
// Это не сложно, но надо ли оно мне сейчас. Пока пусть работает как есть
//
//
// Блин, надо было пакет назвал Оло, или олло, был бы оло.лог(), оло.л()