package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/api"
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

// Check implements api.Operator. проверяет соединение с базой данных, в случае ошибки возвращает error!=nil
func (d *Database) Check(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Get implements api.Operator.
func (*Database) Get(ctx context.Context, req ...metrica.Metrica) (resp []metrica.Metrica, err error) {
	panic("unimplemented")
}

// List implements api.Operator.
func (*Database) List(ctx context.Context) (counterNames []string, gaugeNames []string, err error) {
	panic("unimplemented")
}

// Update implements api.Operator.
func (*Database) Update(ctx context.Context, req ...metrica.Metrica) (resp []metrica.Metrica, err error) {
	panic("unimplemented")
}

func (d Database) Database() *sql.DB {
	return d.db
}

func (d Database) Context() context.Context {
	// todo не понимаю можно ли так передавать контекст по значению
	return d.ctx
}

var _ api.Operator = (*Database)(nil)
