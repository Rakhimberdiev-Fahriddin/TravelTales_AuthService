package handler

import (
	"log/slog"
	"my_module/storage/postgres"
)

type Handler struct {
	User   *postgres.UserRepo
	Logger *slog.Logger
}

func NewHandler(user *postgres.UserRepo, logger *slog.Logger) *Handler {
	return &Handler{
		User:   user,
		Logger: logger,
	}
}
