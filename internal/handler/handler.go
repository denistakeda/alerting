package handler

import (
	s "github.com/denistakeda/alerting/internal/storage"
)

type Handler struct {
	storage s.Storage
	hashKey string
}

func New(storage s.Storage, hashKey string) *Handler {
	return &Handler{
		storage: storage,
		hashKey: hashKey,
	}
}
