package models

import "github.com/google/uuid"

type Subscription struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"Premium"`
	Price     int       `json:"price" example:"100"`
	UserID    uuid.UUID `json:"user_id" example:"11111111-1111-1111-1111-111111111111"`
	StartDate string    `json:"start_date" example:"2025-11"`
	EndDate   string    `json:"end_date,omitempty" example:"2026-11"`
}
