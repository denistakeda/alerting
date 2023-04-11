package main

// @Title Alerting Service API
// @Description Service of metrics and alerting
// @Version 1.0

// @Contact.email denis.takeda@gmail.com

// @Host http://127.0.0.1:8080/

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/denistakeda/alerting/docs"
	servercfg "github.com/denistakeda/alerting/internal/config/server"
	"github.com/denistakeda/alerting/internal/handler"
	"github.com/denistakeda/alerting/internal/services/loggerservice"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/dbstorage"
	"github.com/denistakeda/alerting/internal/storage/filestorage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	conf, err := servercfg.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration: %v", conf)

	logService := loggerservice.New()

	storage, err := getStorage(conf, logService)
	if err != nil {
		log.Fatal(err)
	}

	r := newRouter()
	apiHandler := handler.New(storage, conf.Key, logService)
	apiHandler.RegisterHandlers(r)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	if err := storage.Close(context.Background()); err != nil {
		log.Printf("Unable to properly stop the storage: %v\n", err)
	}
}

func newRouter() *gin.Engine {
	r := gin.New()

	r.RedirectTrailingSlash = false
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Recovery())
	r.Use(logger.SetLogger())

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

func getStorage(conf servercfg.Config, logService *loggerservice.LoggerService) (s.Storage, error) {
	if conf.DatabaseDSN != "" {
		return dbstorage.NewDBStorage(conf.DatabaseDSN, conf.Key, logService)
	} else if conf.StoreFile != "" {
		return filestorage.NewFileStorage(context.Background(), conf.StoreFile, conf.StoreInterval, conf.Restore, conf.Key, logService)
	} else {
		return memstorage.NewMemStorage(conf.Key, logService), nil
	}
}
