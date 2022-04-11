package history

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type History struct {
	storage Storage
}

func New(s Storage) History {
	return History{storage: s}
}

func (h History) Create(m *model.History) error {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()

	if err := h.storage.Create(m); err != nil {
		return fmt.Errorf("h.storage.Create(): %w", err)
	}

	return nil
}

func (h History) ByShortURLID(ID uuid.UUID) (model.Histories, error) {
	return h.storage.ByShortURLID(ID)
}

func (h History) ByShortURLIDAndDates(ID uuid.UUID, from, to time.Time) (model.Histories, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()

	if fromUnix > toUnix {
		fromUnix, toUnix = toUnix, fromUnix
	}

	return h.storage.ByShortURLIDAndDates(ID, fromUnix, toUnix)
}

func (h History) All() (model.Histories, error) {
	return h.storage.All()
}
