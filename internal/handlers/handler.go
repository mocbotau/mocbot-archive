package handlers

import (
	"github.com/mocbotau/mocbot-archive/internal/database"
)

// Handler is the HTTP handler for the API.
type Handler struct {
	db *database.DB
}

// NewHandler creates a new Handler instance.
func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}
