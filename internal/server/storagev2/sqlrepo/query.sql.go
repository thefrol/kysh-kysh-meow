// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: query.sql

package sqlrepo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const counter = `-- name: Counter :one
SELECT
    id,
    delta
FROM
    counters
WHERE
    id=$1
`

func (q *Queries) Counter(ctx context.Context, id string) (Counter, error) {
	row := q.db.QueryRow(ctx, counter, id)
	var i Counter
	err := row.Scan(&i.ID, &i.Delta)
	return i, err
}

const gauge = `-- name: Gauge :one
SELECT 
    id,
    value 
FROM 
    gauges 
WHERE
    id=$1
`

func (q *Queries) Gauge(ctx context.Context, id string) (Gauge, error) {
	row := q.db.QueryRow(ctx, gauge, id)
	var i Gauge
	err := row.Scan(&i.ID, &i.Value)
	return i, err
}

const incrementCounter = `-- name: IncrementCounter :one
INSERT INTO
    counters(id, delta) --todo какая ещё дельта????
VALUES
    ($1,$2)
ON CONFLICT (id)
    DO
    UPDATE SET
        value=$2
RETURNING id, delta
`

type IncrementCounterParams struct {
	ID    string
	Delta pgtype.Int8
}

func (q *Queries) IncrementCounter(ctx context.Context, arg IncrementCounterParams) (Counter, error) {
	row := q.db.QueryRow(ctx, incrementCounter, arg.ID, arg.Delta)
	var i Counter
	err := row.Scan(&i.ID, &i.Delta)
	return i, err
}

const list = `-- name: List :many
SELECT
    'counter',
    id
FROM
    counters
UNION SELECT
    'gauge',
    id
FROM gauges
`

type ListRow struct {
	Column1 string
	ID      string
}

func (q *Queries) List(ctx context.Context) ([]ListRow, error) {
	rows, err := q.db.Query(ctx, list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRow
	for rows.Next() {
		var i ListRow
		if err := rows.Scan(&i.Column1, &i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGauge = `-- name: UpdateGauge :one
INSERT INTO
    gauges(id,value)
VALUES
    ($1,$2)
ON CONFLICT (id)
    DO
    UPDATE SET
        value=$2
RETURNING id, value
`

type UpdateGaugeParams struct {
	ID    string
	Value pgtype.Float8
}

func (q *Queries) UpdateGauge(ctx context.Context, arg UpdateGaugeParams) (Gauge, error) {
	row := q.db.QueryRow(ctx, updateGauge, arg.ID, arg.Value)
	var i Gauge
	err := row.Scan(&i.ID, &i.Value)
	return i, err
}
