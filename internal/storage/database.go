package storage

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/thefrol/kysh-kysh-meow/internal/metrica"
	"github.com/thefrol/kysh-kysh-meow/internal/server/router/httpio"
	"github.com/thefrol/kysh-kysh-meow/lib/retry"
	"github.com/thefrol/kysh-kysh-meow/lib/retry/fails"
)

var (
	ErrorInitDatabase = errors.New("невозможно инициализировать базу данных, создать таблицы counter и gauge")
)

const (
	initQuery = "CREATE TABLE IF NOT EXISTS counters(id TEXT PRIMARY KEY, delta BIGINT);" +
		"CREATE TABLE IF NOT EXISTS gauges(id TEXT PRIMARY KEY, value DOUBLE PRECISION);" // todo вот тут я бы уже делал ошибку обертку, соишком много тонкостей и синтакс и ещё соединения

	queryGetCounter = "SELECT id, delta FROM counters WHERE id=$1;"
	queryGetGauge   = "SELECT id, value FROM gauges WHERE id=$1;"

	queryList = "SELECT 'counter',id FROM counters UNION SELECT 'gauge',id from gauges;"

	queryUpsertCounter = "INSERT INTO counters VALUES ($1,$2) ON CONFLICT (id) DO UPDATE SET delta=counters.delta+$2;"
	queryUpsertGauge   = "INSERT INTO gauges VALUES ($1,$2) ON CONFLICT (id) DO UPDATE SET value=$2;"
)

// Database это приложение сервера, которое умеет работать с базой данных, и другими хранилищами. Со всеми вещами от которого, он зависит.
type Database struct {
	db *sql.DB
}

// New cоздает новый объект приложения, получая на вход параметры конфигурации
func NewDatabase(db *sql.DB) (*Database, error) {
	// инициализуем таблицы для гаужей и каунтеров
	//
	// todo использовать транзации с отменой
	err :=
		retry.This(
			func() error {
				_, err := db.Exec(initQuery)
				return err
			},
			retry.If(fails.OnDial),
			retry.Attempts(3),
			retry.DelaySeconds(1, 3, 5, 7),
			retry.OnRetry(
				func(i int, err error) {
					log.Info().Msgf("Попытка инициализации базы %v: %v", i, err)
				}))

	if err != nil {
		return nil, ErrorInitDatabase
	}
	return &Database{
		db: db,
	}, nil
}

// Check implements api.Operator. проверяет соединение с базой данных, в случае ошибки возвращает error!=nil
func (d *Database) Check(ctx context.Context) error {
	// todo
	//
	// в pgx есть прикольная функция для этого и можно выделить этот метод из Storage
	//
	// конечно, хотчется вынести этот метод в app
	return d.db.PingContext(ctx)
}

// Get implements api.Operator.
func (d *Database) Get(ctx context.Context, req ...metrica.Metrica) (resp []metrica.Metrica, err error) {
	resp = make([]metrica.Metrica, 0, len(req))

	for _, r := range req {
		switch r.MType {
		case "counter":
			result := metrica.Metrica{MType: r.MType, Delta: new(int64)}
			rw := d.db.QueryRowContext(ctx, queryGetCounter, r.ID)
			err := rw.Scan(&result.ID, result.Delta)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, httpio.ErrorNotFoundMetric // todo, вот тут можно упаковку ошибки сделать впринципе
				}
				return nil, err
			}
			resp = append(resp, result)

		case "gauge":
			result := metrica.Metrica{MType: r.MType, Value: new(float64)}
			rw := d.db.QueryRowContext(ctx, queryGetGauge, r.ID)
			err := rw.Scan(&result.ID, result.Value)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, httpio.ErrorNotFoundMetric // todo упаковать ошибку???
				}
				return nil, err
			}
			resp = append(resp, result)
		default:
			return nil, httpio.ErrorUnknownMetricType
		}
	}

	return resp, nil
}

// List implements api.Operator.
func (d *Database) List(ctx context.Context) (counterNames []string, gaugeNames []string, err error) {
	metrics := make(map[string][]string, 2)
	metrics["counter"] = make([]string, 0, 10)
	metrics["gauge"] = make([]string, 0, 30) // todo, если бы мы знали количество возвращаемых строк, то можно было бы тут поменьше памяти выделять

	rs, err := d.db.QueryContext(ctx, queryList)
	if err != nil {
		return nil, nil, err
	}

	var mtype, id string
	for rs.Next() {
		err := rs.Scan(&mtype, &id)
		if err != nil {
			return nil, nil, err
		}
		metrics[mtype] = append(metrics[mtype], id)
	}

	if err := rs.Err(); err != nil {
		return nil, nil, err
	}

	return metrics["counter"], metrics["gauge"], nil
}

// Этот мьютекс добавлен потому что транзакции друг друга
// блокируют, довольно жесткий костыль. От него надо избалсяться
// как можно скорее
var mu = sync.Mutex{}

// Update implements api.Operator.
func (d *Database) Update(ctx context.Context, req ...metrica.Metrica) (resp []metrica.Metrica, err error) {
	mu.Lock()
	defer mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// кажется нужен сквошинг метрик, иначе они друг друууууугу мешают
	counters, gauges := squash(req...)

	for k, v := range counters {

		// if v == nil {
		// 	return nil, api.ErrorDeltaEmpty
		// }

		_, err := tx.ExecContext(ctx, queryUpsertCounter, k, v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range gauges {
		// if r.Value == nil {
		// 	return nil, api.ErrorValueEmpty
		// }

		_, err := tx.ExecContext(ctx, queryUpsertGauge, k, v)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	/// возвращаем обновленные метрики
	return d.Get(ctx, req...)
}

var _ httpio.Operator = (*Database)(nil)

func squash(ms ...metrica.Metrica) (counters map[string]int64, gauges map[string]float64) {
	counters = make(map[string]int64)
	gauges = make(map[string]float64)

	for _, m := range ms {
		switch m.MType {
		case "gauge":
			gauges[m.ID] = *m.Value
		case "counter":
			counters[m.ID] += *m.Delta
		default:
			log.Warn().Msg("Unknown metric type")
		}
	}

	return counters, gauges
}

// TODO
//
// Альтернативный способ получить все метрики одним запросом
//
// SELECT 'counter', id, delta, NULL FROM counters WHERE id IN ('test1','test2')
//   UNION
//   SELECT 'gauge', id, NULL, value FROM gauges WHERE id IN ('gaug2', 'gaug3')
//
// получим три столбца прям как в Metrica
//
//  ?column? |  id   | ?column? | delta
// ----------+-------+----------+-------
//  counter  | test1 |          |    20
//  counter  | test2 |          |    40
//  gauge    | gaug3 |  50.22   |
//
