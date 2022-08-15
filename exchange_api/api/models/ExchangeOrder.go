package models

import (
	"github.com/google/uuid"
	"time"
)

// ExchangeOrder model
type ExchangeOrder struct {
	TransactionId uuid.UUID `json:"transaction_id" gorm:"primary_key"`
	From          string    `json:"from" gorm:"not null;unique;column:fro"`
	Fund          string    `json:"fund"`
	Amount        float64   `json:"amt"`
	Re            string    `json:"re"`
	ReceivedAt    time.Time `json:"receivedAt"`
	ValidatedAt   time.Time `json:"validatedAt"`
	StoredAt      time.Time `json:"storedAt"`
}
