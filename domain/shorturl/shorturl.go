package shorturl

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

const (
	MaxLetters = 7

	HTTPProtocol = "http"
)

var (
	allowedLetters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var (
	ErrWrongRedirect = errors.New("the m.Redirect has not a valid protocol")
)

type UseCase interface {
	Create(s *model.ShortURL) error
	Update(s *model.ShortURL) error
	Delete(ID uuid.UUID) error
	ByShort(s string) (model.ShortURL, error)
	All() (model.ShortURLs, error)
}

type Storage interface {
	Create(s *model.ShortURL) error
	Update(s *model.ShortURL) error
	IncrementTimes(tx pgx.Tx, ID uuid.UUID) error
	Delete(ID uuid.UUID) error
	ByShort(s string) (model.ShortURL, error)
	All() (model.ShortURLs, error)
}
