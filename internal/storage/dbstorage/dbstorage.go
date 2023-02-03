package dbstorage

import (
	"context"
	"database/sql"
	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

type DBStorage struct {
	db      *sqlx.DB
	hashKey string
	logger  zerolog.Logger
}

func New(
	ctx context.Context,
	dsn string,
	hashKey string,
	logService *loggerservice.LoggerService,
) (*DBStorage, error) {
	// This line only required to pass tests for 10th iteration that check the usage of database/sql
	_ = sql.Drivers()

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	if err := bootstrapDatabase(ctx, db); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap database")
	}

	return &DBStorage{
		db:      db,
		hashKey: hashKey,
		logger:  logService.ComponentLogger("DBStorage"),
	}, nil
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

func (dbs *DBStorage) UpdateAll(ctx context.Context, metrics []*metric.Metric) error {
	tx, err := dbs.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start a transaction")
	}

	stmt, err := tx.Prepare(`
		INSERT INTO metrics (id, mtype, value, delta)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id, mtype)
		DO UPDATE SET
		    value = $3,
			delta = metrics.delta + $4
	`)

	if err != nil {
		return errors.Wrap(err, "failed to prepare the update query")
	}

	defer stmt.Close()

	for _, met := range metrics {
		if _, err := stmt.Exec(met.ID, met.MType, met.Value, met.Delta); err != nil {
			if err := tx.Rollback(); err != nil {
				dbs.logger.Fatal().Err(err).Msg("update drivers: unable to rollback")
			}
			return errors.Wrapf(err, "failed to exec query with metric %v", met)
		}
	}

	if err := tx.Commit(); err != nil {
		dbs.logger.Fatal().Err(err).Msg("update drivers: unable to commit")
	}

	return nil
}

// TODO: return error
func (dbs *DBStorage) All(ctx context.Context) []*metric.Metric {
	result := make([]*metric.Metric, 0)

	err := dbs.db.SelectContext(ctx, &result, `
		SELECT *
		FROM metrics
	`)
	if err != nil {
		dbs.logger.Error().Err(err).Msg("failed to query list of all metrics")
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
		    delta BIGINT,
		    UNIQUE (id, mtype)
		);

		CREATE UNIQUE INDEX IF NOT EXISTS id_mtype_index
		ON metrics (id, mtype)
	`)

	if err != nil {
		return errors.Wrap(err, "unable to create table 'metrics'")
	}

	return nil
}
