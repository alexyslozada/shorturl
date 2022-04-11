package history

import (
	"time"

	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type UseCase interface {
	Create(h *model.History) error
	ByShortURLID(ID uuid.UUID) (model.Histories, error)
	ByShortURLIDAndDates(ID uuid.UUID, from, to time.Time) (model.Histories, error)
	All() (model.Histories, error)
}

type Storage interface {
	Create(h *model.History) error
	ByShortURLID(ID uuid.UUID) (model.Histories, error)
	ByShortURLIDAndDates(ID uuid.UUID, from, to int64) (model.Histories, error)
	All() (model.Histories, error)
}
