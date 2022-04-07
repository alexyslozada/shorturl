package model

import "github.com/google/uuid"

type Permission struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CanCreate bool      `json:"can_create"`
	CanUpdate bool      `json:"can_update"`
	CanDelete bool      `json:"can_delete"`
	CanSelect bool      `json:"can_select"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type Permissions []Permission
