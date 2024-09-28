package handlers

import (
	"github.com/go-playground/validator/v10"
	"go-identity/pkg/models"
	"log/slog"
)

type Handler struct {
	Model     models.Model
	Logger    *slog.Logger
	Validator *validator.Validate
}
