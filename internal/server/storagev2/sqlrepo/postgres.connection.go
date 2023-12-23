package sqlrepo

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/thefrol/kysh-kysh-meow/configs/migrations"
)

// StartPostgres соединяется с базой данных постргрес
// при помощи строки подключения cs и активирует миграции,
// возвращает рабочее и проверенное подключение.
func StartPostgres(cs string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", cs)
	if err != nil {
		return nil, fmt.Errorf("postgres.adapter: %w", err)
	}

	goose.SetBaseFS(migrations.EmbedMigrations)
	err = goose.Up(conn, ".")
	if err != nil {
		return nil, fmt.Errorf("postgres.adapter: %w", err)
	}

	return conn, nil

}
