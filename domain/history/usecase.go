package history

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

type History struct {
	storage         Storage
	useCaseShortURL UseCaseShortURL
}

func New(s Storage, ucs UseCaseShortURL) History {
	return History{storage: s, useCaseShortURL: ucs}
}

func (h History) CreateWithTx(tx pgx.Tx, m *model.History) error {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()

	err := h.storage.CreateWithTx(tx, m)
	if err != nil {
		return err
	}

	err = h.useCaseShortURL.IncrementTimes(tx, m.ShortURLID)
	if err != nil {
		return err
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
