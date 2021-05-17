package persistence

import (
	_ "embed" // sql query embedded
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // db driver
	"github.com/jmoiron/sqlx"

	"github.com/popsu/go-website-watcher/internal/model"
)

//go:embed sql/insert.sql
var sqlInsertQuery string

type PostgresStore struct {
	db     *sqlx.DB
	logger *log.Logger
}

func NewPostgresStore(dburl string, logger *log.Logger) (*PostgresStore, error) {
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

	return &PostgresStore{
		db:     db,
		logger: logger,
	}, nil
}

func (ps *PostgresStore) InsertMessage(msg *model.Message) error {
	ct, err := ps.db.Exec(sqlInsertQuery,
		msg.ID,
		msg.CreatedAt,
		msg.URL,
		msg.Error,
		msg.RegexpPattern,
		msg.RegexpMatch,
		msg.StatusCode,
		msg.TimeToFirstByte,
	)

	if err != nil {
		return err
	}

	n, err := ct.RowsAffected()
	if err != nil {
		return err
	}

	ps.logger.Printf("Inserted %d rows successfully", n)

	return nil
}

// Close closes the database connection
func (ps *PostgresStore) Close() error {
	if ps.db != nil {
		err := ps.db.Close()
		if err != nil {
			return err
		}

		ps.logger.Println("DB connection closed")
	}

	return nil
}
