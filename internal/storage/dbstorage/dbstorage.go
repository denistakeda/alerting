package dbstorage

import (
	"context"
	"database/sql"
	"github.com/denistakeda/alerting/internal/metric"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"log"
	"time"
)

type DBStorage struct {
	db      *sql.DB
	hashKey string
}

func New(dsn string, hashKey string) (*DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	if err := bootstrapDatabase(db); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap database")
	}

	return &DBStorage{db: db, hashKey: hashKey}, nil
}

func (dbs *DBStorage) Get(metricType metric.Type, metricName string) (*metric.Metric, bool) {
	var met metric.Metric
	err := dbs.db.QueryRow(`
		SELECT id, mtype, value, delta
		FROM metrics
		WHERE id=$1 AND mtype=$2
	`, metricName, metricType).
		Scan(&met.ID, &met.MType, &met.Value, &met.Delta)
	if err != nil {
		return nil, false
	}
	met.FillHash(dbs.hashKey)
	return &met, true
}

func (dbs *DBStorage) Update(met *metric.Metric) (*metric.Metric, error) {
	oldMet, ok := dbs.Get(met.Type(), met.Name())
	newMet := metric.Update(oldMet, met)
	var row *sql.Row
	if ok {
		row = dbs.db.QueryRow(`
			UPDATE metrics
			SET value = $1,
				delta = $2
			WHERE id = $3 AND mtype = $4
		`, newMet.Value, newMet.Delta, newMet.ID, newMet.MType)
	} else {
		row = dbs.db.QueryRow(`
			INSERT INTO metrics (id, mtype, value, delta)
			VALUES ($1, $2, $3, $4)
		`, newMet.ID, newMet.MType, newMet.Value, newMet.Delta)
	}
	if row.Err() != nil {
		return nil, errors.Wrap(row.Err(), "unable to update metric")
	}
	newMet.FillHash(dbs.hashKey)
	return newMet, nil
}

// TODO: return error
func (dbs *DBStorage) All() []*metric.Metric {
	result := make([]*metric.Metric, 0)

	rows, err := dbs.db.Query(`
		SELECT id, mtype, value, delta
		FROM metrics
	`)
	if err != nil {
		log.Println("failed to query list of all metrics")
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var m metric.Metric
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			log.Println("failed to parse metric")
			continue
		}
		m.FillHash(dbs.hashKey)
		result = append(result, &m)
	}

	err = rows.Err()
	if err != nil {
		log.Println("error while iterate over list of metrics")
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

func (dbs *DBStorage) Close() error {
	return dbs.db.Close()
}

func bootstrapDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
    		id VARCHAR(256),
		    mtype VARCHAR(10),
		    value NUMERIC,
		    delta BIGINT
		)
	`)
	if err != nil {
		return errors.Wrap(err, "unable to create table 'metrics'")
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS id_mtype_index
		ON metrics (id, mtype)
	`)
	if err != nil {
		return errors.Wrap(err, "unable to create index for table 'metrics'")
	}

	return nil
}
