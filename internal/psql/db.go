package psql

import (
	"time"

	"github.com/jmoiron/sqlx"
	// db driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

func New(dburl string) (*sqlx.DB, error) {
	println(dburl)

	db, err := sqlx.Open("pgx", dburl)
	if err != nil {
		return nil, err
	}

	// values from https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
