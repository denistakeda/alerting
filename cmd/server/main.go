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
	r.LoadHTMLGlob("cmd/server/templates/*")
	r.Run()
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
