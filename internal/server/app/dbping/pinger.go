// Этот юзкейс пингует базу данных. Я хотел
// чтобы это было как-то совсем отдельно работало.
// то есть, по сути, чтобы пинговалась база, даже
// если мы работаем в памяти. Ну такое задание что сказать.
//
// Есть в этом немного абсурда, но для практики это прям классный
// вариант
package dbping

import (
	"context"
	"database/sql"
)

type Pinger struct {
	Connection *sql.DB
}

func (p *Pinger) Ping(ctx context.Context) error {
	err := p.Connection.PingContext(ctx)
	return err
}
