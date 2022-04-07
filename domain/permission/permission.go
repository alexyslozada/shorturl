package permission

import (
	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type UseCase interface {
	Create(p *model.Permission) error
	Update(p *model.Permission) error
	Delete(ID uuid.UUID) error
	ByUserID(ID uuid.UUID) (model.Permission, error)
	All() (model.Permissions, error)
}

type Storage interface {
	Create(p *model.Permission) error
	Update(p *model.Permission) error
	Delete(ID uuid.UUID) error
	ByUserID(ID uuid.UUID) (model.Permission, error)
	All() (model.Permissions, error)
}
