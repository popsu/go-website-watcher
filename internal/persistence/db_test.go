package persistence

import (
	"errors"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // db driver
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/popsu/go-website-watcher/internal/model"
	"github.com/stretchr/testify/assert"
)

const (
	migrationVersion = 1
	migrationFolder  = "file:../../sql/migrations"
	dockerImage      = "postgres"
	dockerTag        = "12.7"
)

// startDB starts a database using dockertest and returns the database url
func startDB(t *testing.T) *url.URL {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to docker: %v", err)
	}

	dburl := "postgres://postgres:pass@localhost:5432/testdb"
	pgURL, err := url.Parse(dburl)
	if err != nil {
		t.Fatalf("Bad postgres URL: %v", err)
	}

	pass, _ := pgURL.User.Password()

	resource, err := pool.Run(dockerImage, dockerTag, []string{
		"POSTGRES_USER=" + pgURL.User.Username(),
		"POSTGRES_DB=" + pgURL.Path[1:],
		"POSTGRES_PASSWORD=" + pass,
	})
	if err != nil {
		t.Fatalf("Could not start postgres container: %v", err)
	}

	t.Cleanup(func() {
		err = pool.Purge(resource)
		if err != nil {
			t.Fatalf("Could not purge container: %v", err)
		}
	})

	// "Adjust" host:port to whatever docker is exposing
	pgURL.Host = resource.GetHostPort("5432/tcp")

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() (err error) {
		db, err := sqlx.Open("pgx", pgURL.String())
		if err != nil {
			return err
		}
		defer func() {
			cerr := db.Close()
			if err == nil {
				err = cerr
			}
		}()

		return db.Ping()
	})

	if err != nil {
		t.Fatalf("Could not connect to postgres container: %v", err)
	}

	return pgURL
}

func migrateUp(db *sqlx.DB, dburl string) error {
	// go-migrations require sql.DB rather than sqlx.DB
	sqlDB := db.DB

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	// Tables
	m, err := migrate.NewWithDatabaseInstance(migrationFolder, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Migrate(migrationVersion)
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func TestInsertMessage(t *testing.T) {
	t.Parallel()

	logger := log.Default()

	dburl := startDB(t).String()

	store, err := NewPostgresStore(dburl, logger)
	if err != nil {
		t.Fatalf("Error connecting to db: %s", err)
	}

	t.Cleanup(func() {
		err = store.Close()
		if err != nil {
			t.Errorf("Error closing the db: %s", err)
		}
	})

	err = migrateUp(store.db, dburl)
	if err != nil {
		t.Fatalf("Error migrating the db: %s", err)
	}

	t.Run("Insert message", func(t *testing.T) {
		t.Parallel()

		id, err := uuid.NewV4()
		if err != nil {
			t.Fatalf("Error creating uuid: %s", err)
		}

		ttfb := (time.Duration)(200)
		createdAt := time.Now()
		url := "http://www.example.com"
		rePattern := "testre123"
		reMatch := true
		statusCode := int32(200)

		testMessage := &model.Message{
			ID:              &id,
			CreatedAt:       &createdAt,
			URL:             &url,
			Error:           nil,
			RegexpPattern:   &rePattern,
			RegexpMatch:     &reMatch,
			StatusCode:      &statusCode,
			TimeToFirstByte: &ttfb,
		}

		err = store.InsertMessage(testMessage)

		assert.NoError(t, err)

		// store.InsertMessage(msg * model.Message)
	})
}
