package dbstorage

import (
	"context"
	"database/sql"
	"github.com/denistakeda/alerting/internal/metric"
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

	if err := bootstrapDatabase(db); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap database")
	}

	return &DBStorage{db: db}, nil
}

func (dbs *DBStorage) Get(metricType metric.Type, metricName string) (*metric.Metric, bool) {
	//TODO implement me
	panic("implement me")
}

func (dbs *DBStorage) Update(metric *metric.Metric) (*metric.Metric, error) {
	//TODO implement me
	panic("implement me")
}

func (dbs *DBStorage) All() []*metric.Metric {
	//TODO implement me
	panic("implement me")
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

func bootstrapDatabase(db *sql.DB) error {
	row := db.QueryRow(`
		CREATE TABLE IF NOT EXISTS metrics (
    		id VARCHAR(256),
		    mtype VARCHAR(10),
		    value NUMERIC(64),
		    delta INT
		)
	`)
	return row.Err()
}
