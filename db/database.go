package db

import "github.com/jmoiron/sqlx"

type IDatabase interface {
    Connect() *sqlx.DB
    Close() error
}
