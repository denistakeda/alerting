package main

import (
	"github.com/denistakeda/alerting/cmd/server/internal/config"
	"github.com/denistakeda/alerting/cmd/server/internal/handler"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/filestorage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration: %v", conf)

	storage := getStorage(conf)

	r := setupRouter(storage)
	r.LoadHTMLGlob("cmd/server/internal/templates/*")
	log.Fatal(r.Run(conf.Address))
}

func setupRouter(storage s.Storage) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false

	h := handler.New(storage)

	r.POST("/update/", h.UpdateMetricHandler2)
	r.POST("/update/:metric_type/:metric_name/:metric_value", h.UpdateMetricHandler)
	r.POST("/value/", h.GetMetricHandler2)
	r.GET("/value/:metric_type/:metric_name", h.GetMetricHandler)
	r.GET("/", h.MainPageHandler)
	return r
}

func getStorage(conf config.Config) s.Storage {
	if conf.StoreFile == "" {
		return memstorage.New()
	} else {
		storage, err := filestorage.New(conf.StoreFile, conf.StoreInterval, conf.Restore)
		if err != nil {
			log.Fatal(err)
		}
		return storage
	}
}
