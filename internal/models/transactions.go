package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID         uuid.UUID `json:"id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	FromUser   string    `json:"fromUser"`
	ToUserID   uuid.UUID `json:"to_user_id"`
	ToUser     string    `json:"toUser"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time
}
