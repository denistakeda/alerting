package main

import (
	"github.com/denistakeda/alerting/cmd/server/handler"
	"github.com/denistakeda/alerting/internal/memstorage"
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	storage := memstorage.NewMemStorage()
	r := setupRouter(storage)
	r.Run()
}

func setupRouter(storage s.Storage) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = false

	r.POST("/update/:metric_type/:metric_name/:metric_value", handler.UpdateMetricHandler(storage))
	return r
}
