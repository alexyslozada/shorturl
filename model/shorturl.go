package model

import "github.com/google/uuid"

type ShortURL struct {
	ID          uuid.UUID `json:"id"`
	IsRandom    bool      `json:"is_random"`
	Short       string    `json:"short"`
	RedirectTo  string    `json:"redirect_to"`
	Description string    `json:"description"`
	Times       int       `json:"times"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

type ShortURLs []ShortURL
