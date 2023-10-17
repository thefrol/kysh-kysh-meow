// Некоторые ошибки временны, и функции их получившие можно просто
// попытаться выполнить ещё раз, и в следующий раз все нормально
// завершиться
//
// В этом пакете собраны функции-проверялки таких ошибок, из
// самых растространенных случаях, таких как http и бд
package fails

import (
	"errors"
	"net"
)

// OnDial возвращает true, если err связана с ошибкой подключения
// то есть ошибка net.OpError, где operr.Op=="dial"
//
// retry.This(func ()error{}, retry.If(fails.OnDial))
func OnDial(err error) bool {
	var oe *net.OpError
	if errors.As(err, &oe) {
		return oe.Op == "dial" // если ошибка в операции dial
	}
	return false
}
