// Здесь мы прорабатываем зависимости приложения. Пакет содержит
// класс App, который умеет работать со своей базой данных, хранилищем и прочим.
package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// App это приложение сервера, которое умеет работать с базой данных, и другими хранилищами. Со всеми вещами от которого, он зависит.
type App struct {
	db *sql.DB
}

const sqldriver = "pgx"

// New cоздает новый объект приложения, получая на вход параметры конфигурации
func New(connString string) (*App, error) {

	db, err := createDataBase(sqldriver, connString)
	if err != nil {
		return nil, fmt.Errorf("не могу создать базу данных %v", err)
	}
	return &App{
		db: db,
	}, nil
}

// CheckConnection проверяет соединение с базой данных, в случае ошибки возвращает error!=nil
func (app App) CheckConnection(ctx context.Context) error {
	return app.db.PingContext(ctx)
}

func (app App) Database() *sql.DB {
	return app.db
}

func createDataBase(driver string, connstring string) (*sql.DB, error) {
	conn, err := sql.Open(driver, connstring)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
