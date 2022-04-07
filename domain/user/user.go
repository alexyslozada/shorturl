package user

import (
	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type UseCase interface {
	Create(u *model.User) error
	Delete(ID uuid.UUID) error
	ByEmail(email string) (model.User, error)
	All() (model.Users, error)
}

type Storage interface {
	Create(u *model.User) error
	Delete(ID uuid.UUID) error
	ByEmail(email string) (model.User, error)
	All() (model.Users, error)
}
