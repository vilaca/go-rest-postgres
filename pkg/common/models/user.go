package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"username"`
	Password string    `json:"password"`
	Enabled  bool      `json:"enabled"`
}
