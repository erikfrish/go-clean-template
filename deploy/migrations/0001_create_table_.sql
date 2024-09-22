-- +goose Up
CREATE SCHEMA IF NOT EXISTS schema_;
CREATE TABLE if not exists schema_.table_ (
    column_ TEXT NOT NULL, 
    PRIMARY KEY (column_)
);
CREATE INDEX if not exists column_pkey ON schema_.table_ (column_);


-- +goose Down
--DROP TABLE schema_.table_;
