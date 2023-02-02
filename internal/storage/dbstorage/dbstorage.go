package dbstorage

import (
	"context"
	"github.com/denistakeda/alerting/internal/metric"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

type DBStorage struct {
	db      *sqlx.DB
	hashKey string
}

func New(ctx context.Context, dsn string, hashKey string) (*DBStorage, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	if err := bootstrapDatabase(ctx, db); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap database")
	}

	return &DBStorage{db: db, hashKey: hashKey}, nil
}

func (dbs *DBStorage) Get(ctx context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool) {
	var met metric.Metric
	err := dbs.db.GetContext(ctx, &met, `
		SELECT *
		FROM metrics
		WHERE id=$1 AND mtype=$2
	`, metricName, metricType)

	if err != nil {
		return nil, false
	}

	met.FillHash(dbs.hashKey)

	return &met, true
}

func (dbs *DBStorage) Update(ctx context.Context, met *metric.Metric) (*metric.Metric, error) {
	oldMet, ok := dbs.Get(ctx, met.Type(), met.Name())
	newMet := metric.Update(oldMet, met)
	var err error
	if ok {
		_, err = dbs.db.NamedExecContext(ctx, `
			UPDATE metrics
			SET value = :value,
				delta = :delta 
			WHERE id = :id AND mtype = :mtype 
		`, newMet)
	} else {
		_, err = dbs.db.NamedExecContext(ctx, `
			INSERT INTO metrics (id, mtype, value, delta)
			VALUES (:id, :mtype, :value, :delta)
		`, newMet)
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to update metric")
	}

	newMet.FillHash(dbs.hashKey)

	return newMet, nil
}

// TODO: return error
func (dbs *DBStorage) All(ctx context.Context) []*metric.Metric {
	result := make([]*metric.Metric, 0)

	err := dbs.db.SelectContext(ctx, &result, `
		SELECT *
		FROM metrics
	`)
	if err != nil {
		log.Println("failed to query list of all metrics")
		return result
	}

	return result
}

func (dbs *DBStorage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := dbs.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) Close(_ context.Context) error {
	return dbs.db.Close()
}

func bootstrapDatabase(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS metrics (
    		id VARCHAR(256),
		    mtype VARCHAR(10),
		    value NUMERIC,
		    delta BIGINT
		);

		CREATE UNIQUE INDEX IF NOT EXISTS id_mtype_index
		ON metrics (id, mtype)
	`)

	if err != nil {
		return errors.Wrap(err, "unable to create table 'metrics'")
	}

	return nil
}
