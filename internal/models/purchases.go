package models

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	ID        uuid.UUID `json:"id"`
	UserID    uint      `json:"user_id"`
	ItemName  string    `json:"item_name"`
	Quantity  uint      `json:"quantity"`
	CreatedAt time.Time
}
