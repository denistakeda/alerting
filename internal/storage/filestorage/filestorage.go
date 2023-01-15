package filestorage

import (
	"encoding/json"
	"github.com/denistakeda/alerting/internal/metric"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"time"
)

type Filestorage struct {
	mstorage *memstorage.Memstorage

	storeFile   string
	storeTicker *time.Ticker
}

func New(storeFile string, storeInterval time.Duration, restore bool) (*Filestorage, error) {
	instance := &Filestorage{
		mstorage:  memstorage.New(),
		storeFile: storeFile,
	}

	if restore {
		if err := instance.restore(); err != nil {
			return nil, errors.Wrap(err, "unable to initiate a Filestorage")
		}
	}

	if storeInterval != 0 {
		instance.storeTicker = time.NewTicker(storeInterval)
		go func() {
			for range instance.storeTicker.C {
				instance.dump()
			}
		}()
	}

	return instance, nil
}

func (fs *Filestorage) Get(metricType metric.Type, metricName string) (*metric.Metric, bool) {
	return fs.mstorage.Get(metricType, metricName)
}

func (fs *Filestorage) Update(updatedMetric *metric.Metric) (*metric.Metric, error) {
	res, err := fs.mstorage.Update(updatedMetric)
	if fs.storeTicker == nil {
		fs.dump()
	}
	return res, err
}

func (fs *Filestorage) All() []*metric.Metric {
	return fs.mstorage.All()
}

func (fs *Filestorage) Close() error {
	fs.storeTicker.Stop()
	fs.dump()
	return fs.mstorage.Close()
}

func (fs *Filestorage) dump() {
	logPrefix := "Filestorage: failed to dump data"
	ms := fs.All()
	if len(ms) == 0 {
		return
	}

	file, err := os.OpenFile(fs.storeFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("%s: failed to open the file \"%s\"", logPrefix, fs.storeFile)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file \"%s\"", fs.storeFile)
		}
	}()
	encoder := json.NewEncoder(file)

	for _, met := range ms {
		if err := encoder.Encode(met); err != nil {
			log.Printf("%s: failed to write metric %v to file %s: %v", logPrefix, met, fs.storeFile, err)
		}
	}
}

func (fs *Filestorage) restore() error {
	file, err := os.OpenFile(fs.storeFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.Wrapf(err, "Failestorage: failed to restore data from file %s", fs.storeFile)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file \"%s\"", fs.storeFile)
		}
	}()
	decoder := json.NewDecoder(file)
	var m metric.Metric
	for {
		err := decoder.Decode(&m)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "Failestorage: failed to restore data from file %s", fs.storeFile)
		}
		if _, err := fs.mstorage.Update(&m); err != nil {
			return errors.Wrapf(err, "Failestorage: failed to restore data from file %s", fs.storeFile)
		}
	}
}
