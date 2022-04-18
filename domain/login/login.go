package login

import "github.com/alexyslozada/shorturl/model"

type UseCase interface {
	Login(email, password string) (string, error)
}

type Storage interface {
	ByEmail(email string) (model.User, error)
}
