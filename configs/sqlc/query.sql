-- name: Counter :one
SELECT
    id,
    delta
FROM
    counters
WHERE
    id=$1;

-- name: Gauge :one
SELECT 
    id,
    value 
FROM 
    gauges 
WHERE
    id=$1;

-- name: List :many
SELECT
    'counter',
    id
FROM
    counters
UNION SELECT
    'gauge',
    id
FROM gauges;

-- name: IncrementCounter :one
INSERT INTO
    counters(id, delta) --todo какая ещё дельта????
VALUES
    ($1,$2)
ON CONFLICT (id)
    DO
    UPDATE SET
        value=$2
RETURNING *;    


-- name: UpdateGauge :one
INSERT INTO
    gauges(id,value)
VALUES
    ($1,$2)
ON CONFLICT (id)
    DO
    UPDATE SET
        value=$2
RETURNING *;
