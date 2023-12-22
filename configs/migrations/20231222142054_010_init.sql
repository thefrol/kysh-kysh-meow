-- +goose Up
-- +goose StatementBegin

CREATE TABLE 
	IF NOT EXISTS 
	counters(
		id TEXT PRIMARY KEY,
		delta BIGINT);

CREATE TABLE
    IF NOT EXISTS
	gauges(
		id TEXT PRIMARY KEY,
		value DOUBLE PRECISION);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE
    gauges;

DROP TABLE
    counters;
-- +goose StatementEnd
