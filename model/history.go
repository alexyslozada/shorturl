package model

import "github.com/google/uuid"

type History struct {
	ID         uuid.UUID `json:"id"`
	ShortURLID uuid.UUID `json:"short_url_id"`
	CreatedAt  int64     `json:"created_at"`
	UpdatedAt  int64     `json:"updated_at"`
}

type Histories []History
