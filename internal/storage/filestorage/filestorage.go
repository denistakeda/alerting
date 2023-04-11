package filestorage

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
)

// Filestorage is an implementation of Storage which stores data in a file.
type Filestorage struct {
	mstorage *memstorage.Memstorage

	storeFile   string
	storeTicker *time.Ticker

	logger zerolog.Logger
}

// NewFileStorage instantiates a new instance of Filestorage
func NewFileStorage(
	ctx context.Context,
	storeFile string,
	storeInterval time.Duration,
	restore bool,
	hashKey string,
	logService *loggerservice.LoggerService,
) (*Filestorage, error) {
	instance := &Filestorage{
		mstorage:  memstorage.NewMemStorage(hashKey, logService),
		storeFile: storeFile,
		logger:    logService.ComponentLogger("Filestorage"),
	}

	if restore {
		if err := instance.restore(ctx); err != nil {
			return nil, errors.Wrap(err, "unable to initiate a Filestorage")
		}
	}

	if storeInterval != 0 {
		instance.storeTicker = time.NewTicker(storeInterval)
		go func() {
			for range instance.storeTicker.C {
				instance.dump(ctx)
			}
		}()
	}

	return instance, nil
}

// Get returns a metric if exists.
func (fs *Filestorage) Get(ctx context.Context, metricType metric.Type, metricName string) (*metric.Metric, bool) {
	return fs.mstorage.Get(ctx, metricType, metricName)
}

// Update updates a metric if exists.
func (fs *Filestorage) Update(ctx context.Context, updatedMetric *metric.Metric) (*metric.Metric, error) {
	res, err := fs.mstorage.Update(ctx, updatedMetric)
	if fs.storeTicker == nil {
		fs.dump(ctx)
	}
	return res, err
}

// UpdateAll updates all the metrics in list.
func (fs *Filestorage) UpdateAll(ctx context.Context, metrics []*metric.Metric) error {
	for _, met := range metrics {
		_, err := fs.Update(ctx, met)
		if err != nil {
			return err
		}
	}

	return nil
}

// All returns all the metrics.
func (fs *Filestorage) All(ctx context.Context) []*metric.Metric {
	return fs.mstorage.All(ctx)
}

// Close closes the connection to db.
func (fs *Filestorage) Close(ctx context.Context) error {
	fs.storeTicker.Stop()
	fs.dump(ctx)
	return fs.mstorage.Close(ctx)
}

func (fs *Filestorage) Ping(_ context.Context) error {
	// For file storage there is no need to do anything on ping
	return nil
}

func (fs *Filestorage) dump(ctx context.Context) {
	logPrefix := "Filestorage: failed to dump data"
	ms := fs.All(ctx)
	if len(ms) == 0 {
		return
	}

	file, err := os.OpenFile(fs.storeFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fs.logger.Error().Err(err).
			Msgf("%s: failed to open the file \"%s\"", logPrefix, fs.storeFile)

		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fs.logger.Error().Err(err).
				Msgf("failed to close file \"%s\"", fs.storeFile)
		}
	}()
	encoder := json.NewEncoder(file)

	for _, met := range ms {
		if err := encoder.Encode(met); err != nil {
			fs.logger.Error().Err(err).
				Msgf("%s: failed to write metric %v to file %s: %v", logPrefix, met, fs.storeFile, err)
		}
	}
}

func (fs *Filestorage) restore(ctx context.Context) error {
	file, err := os.OpenFile(fs.storeFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.Wrapf(err, "Failestorage: failed to restore data from file %s", fs.storeFile)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fs.logger.Error().Err(err).
				Msgf("failed to close file \"%s\"", fs.storeFile)
		}
	}()
	decoder := json.NewDecoder(file)
	for {
		var m metric.Metric
		err := decoder.Decode(&m)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "Failestorage: failed to restore data from file %s", fs.storeFile)
		}
		fs.mstorage.Replace(ctx, &m)
	}
}
