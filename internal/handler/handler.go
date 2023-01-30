package handler

import (
	s "github.com/denistakeda/alerting/internal/storage"
	"github.com/denistakeda/alerting/internal/storage/dbstorage"
)

type Handler struct {
	storage   s.Storage
	dbStorage *dbstorage.DBStorage
	hashKey   string
}

func New(storage s.Storage, dbStorage *dbstorage.DBStorage, hashKey string) *Handler {
	return &Handler{
		storage:   storage,
		hashKey:   hashKey,
		dbStorage: dbStorage,
	}
}
