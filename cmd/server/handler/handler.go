package handler

import (
	s "github.com/denistakeda/alerting/internal/storage"
)

type handler struct {
	storage s.Storage
}

func New(storage s.Storage) *handler {
	return &handler{
		storage: storage,
	}
}
