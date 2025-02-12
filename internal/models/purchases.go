package models

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ItemName  string    `json:"item_name"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time
}
