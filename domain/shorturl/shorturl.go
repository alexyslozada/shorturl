package shorturl

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

const (
	MaxLetters = 7
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

type Logger interface {
	Errorw(msg string, keysAndValues ...interface{})
}

type UseCaseHistory interface {
	CreateWithTx(tx pgx.Tx, m *model.History) error
}

type UseCaseDB interface {
	Tx() (pgx.Tx, error)
}
