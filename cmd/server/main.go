package main

import (
	"github.com/denistakeda/alerting/cmd/server/handler"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/memstorage"
	"github.com/gin-gonic/gin"
)

func main() {
	storage := memstorage.New()
	r := setupRouter(storage)
	r.Run()
}

func setupRouter(storage s.Storage) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.LoadHTMLGlob("cmd/server/templates/*")

	r.POST("/update/:metric_type/:metric_name/:metric_value", handler.UpdateMetricHandler(storage))
	r.GET("/value/:metric_type/:metric_name", handler.GetMetricHandler(storage))
	r.GET("/", handler.MainPageHandler(storage))
	return r
}
