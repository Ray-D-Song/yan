package infra

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(config *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(config.DB.Driver, config.DB.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}
