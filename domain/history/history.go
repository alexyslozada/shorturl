package history

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

type UseCase interface {
	CreateWithTx(tx pgx.Tx, h *model.History) error
	ByShortURLID(ID uuid.UUID) (model.Histories, error)
	ByShortURLIDAndDates(ID uuid.UUID, from, to time.Time) (model.Histories, error)
	All() (model.Histories, error)
}

type Storage interface {
	CreateWithTx(tx pgx.Tx, h *model.History) error
	ByShortURLID(ID uuid.UUID) (model.Histories, error)
	ByShortURLIDAndDates(ID uuid.UUID, from, to int64) (model.Histories, error)
	All() (model.Histories, error)
}

type UseCaseShortURL interface {
	IncrementTimes(tx pgx.Tx, ID uuid.UUID) error
}
