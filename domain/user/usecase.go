package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexyslozada/shorturl/model"
)

type User struct {
	storage Storage
}

func New(s Storage) User {
	return User{storage: s}
}

func (u User) Create(m *model.User) error {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hash)
	m.Active = true

	return u.storage.Create(m)
}

func (u User) Delete(ID uuid.UUID) error {
	return u.storage.Delete(ID)
}

func (u User) ByEmail(email string) (model.User, error) {
	return u.storage.ByEmail(email)
}

func (u User) All() (model.Users, error) {
	return u.storage.All()
}
