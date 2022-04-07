package model

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FullName  string    `json:"full_name"`
	Active    bool      `json:"active"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type Users []User
