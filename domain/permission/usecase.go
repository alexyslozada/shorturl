package permission

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type Permission struct {
	storage Storage
}

func New(s Storage) Permission {
	return Permission{storage: s}
}

func (p Permission) Create(m *model.Permission) error {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()

	return p.storage.Create(m)
}

func (p Permission) Update(m *model.Permission) error {
	m.UpdatedAt = time.Now().Unix()

	return p.storage.Update(m)
}

func (p Permission) Delete(ID uuid.UUID) error {
	return p.storage.Delete(ID)
}

func (p Permission) ByUserID(ID uuid.UUID) (model.Permission, error) {
	return p.storage.ByUserID(ID)
}

func (p Permission) All() (model.Permissions, error) {
	return p.storage.All()
}

func (p Permission) HasPermission(userID uuid.UUID, method string) (bool, error) {
	permission, err := p.ByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("%s %w", "permission.HasPermission()", err)
	}

	switch method {
	case http.MethodPost:
		return permission.CanCreate, nil
	case http.MethodGet:
		return permission.CanSelect, nil
	case http.MethodPut:
		return permission.CanUpdate, nil
	case http.MethodDelete:
		return permission.CanDelete, nil
	}

	return false, model.ErrMethodNotAllowed
}
