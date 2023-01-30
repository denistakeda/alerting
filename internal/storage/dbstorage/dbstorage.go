package dbstorage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func New(dsn string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}
	return &DBStorage{db: db}, nil
}

func (dbs *DBStorage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := dbs.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) Close() error {
	return dbs.db.Close()
}
