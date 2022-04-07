package shorturl

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
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
