package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database это приложение сервера, которое умеет работать с базой данных, и другими хранилищами. Со всеми вещами от которого, он зависит.
type Database struct {
	db  *sql.DB
	ctx context.Context
}

const sqldriver = "pgx"

// New cоздает новый объект приложения, получая на вход параметры конфигурации
func NewPostGresDatabase(ctx context.Context, connString string) (*Database, error) {

	db, err := sql.Open(sqldriver, connString)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:  db,
		ctx: ctx,
	}, nil
}

// CheckConnection проверяет соединение с базой данных, в случае ошибки возвращает error!=nil
func (d Database) CheckConnection(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d Database) Database() *sql.DB {
	return d.db
}

func (d Database) Context() context.Context {
	// todo не понимаю можно ли так передавать контекст по значению
	return d.ctx
}
