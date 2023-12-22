-- +goose Up
-- +goose StatementBegin

CREATE TABLE 
	IF NOT EXISTS 
	counters(
		id TEXT PRIMARY KEY,
		value BIGINT NOT NULL);

CREATE TABLE
    IF NOT EXISTS
	gauges(
		id TEXT PRIMARY KEY,
		value DOUBLE PRECISION NOT NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE
    gauges;

DROP TABLE
    counters;
-- +goose StatementEnd
