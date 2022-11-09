package shorturl

import (
	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

const (
	MaxLetters = 7

	HTTPProtocol = "http"
)

var (
	allowedLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type UseCase interface {
	Create(s *model.ShortURL) error
	Update(s *model.ShortURL) error
	Delete(ID uuid.UUID) error
	ByShort(shortURL string) (model.ShortURL, error)
	ByShortToRedirect(s string) (model.ShortURL, error)
	All(limit, offset int) (model.ShortURLs, error)
}

type Storage interface {
	Create(s *model.ShortURL) error
	Update(s *model.ShortURL) error
	IncrementTimes(ID uuid.UUID) error
	Delete(ID uuid.UUID) error
	ByShort(s string) (model.ShortURL, error)
	All(limit, offset int) (model.ShortURLs, error)
}

type UseCaseHistory interface {
	Create(m *model.History) error
}
