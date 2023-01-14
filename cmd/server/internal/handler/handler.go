package handler

import (
	s "github.com/denistakeda/alerting/internal/storage"
)

type Handler struct {
	storage s.Storage
}

func New(storage s.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
