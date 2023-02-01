package main

import (
	"github.com/denistakeda/alerting/internal/config/server"
	"github.com/denistakeda/alerting/internal/handler"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/dbstorage"
	"github.com/denistakeda/alerting/internal/storage/filestorage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := servercfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration: %v", conf)

	storage, err := getStorage(conf)
	if err != nil {
		log.Fatal(err)
	}

	r := setupRouter(storage, conf.Key)
	r.LoadHTMLGlob("internal/templates/*")
	serverChan := runServer(r, conf.Address)
	interruptChan := handleInterrupt()
	select {
	case serverError := <-serverChan:
		log.Println(serverError)
	case <-interruptChan:
		log.Println("Program was interrupted")
	}

	stopServer(storage)
}

func stopServer(storage s.Storage) {
	if err := storage.Close(); err != nil {
		log.Printf("Unable to properly stop the storage: %v\n", err)
	}
}

func setupRouter(storage s.Storage, hashKey string) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	h := handler.New(storage, hashKey)

	r.POST("/update/", h.UpdateMetricHandler2)
	r.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetricHandler)
	r.POST("/value/", h.GetMetricHandler2)
	r.GET("/value/:metric_type/:metric_name", h.GetMetricHandler)
	r.GET("/ping", h.PingHandler)
	r.GET("/", h.MainPageHandler)
	return r
}

func runServer(r *gin.Engine, address string) <-chan error {
	out := make(chan error)
	go func() {
		err := r.Run(address)
		out <- err
	}()
	return out
}

func handleInterrupt() <-chan os.Signal {
	out := make(chan os.Signal, 2)
	signal.Notify(out, os.Interrupt)
	signal.Notify(out, syscall.SIGTERM)
	return out
}

func getStorage(conf servercfg.Config) (s.Storage, error) {
	if conf.DatabaseDSN != "" {
		return dbstorage.New(conf.DatabaseDSN, conf.Key)
	} else if conf.StoreFile != "" {
		return filestorage.New(conf.StoreFile, conf.StoreInterval, conf.Restore, conf.Key)
	} else {
		return memstorage.New(conf.Key), nil
	}
}
