package main

import (
	"log"
	"net/http"

	"github.com/denistakeda/alerting/cmd/server/handler"
	"github.com/denistakeda/alerting/internal/memstorage"
)

func main() {
	storage := memstorage.NewMemStorage()
	http.HandleFunc("/update/", handler.UpdateHandler(storage))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
