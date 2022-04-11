package history

import (
	"context"
	"fmt"
	"time"

	"github.com/alexyslozada/shorturl/model"

	"github.com/google/uuid"
)

type History struct {
	storage  Storage
	shortURL UseCaseShortURL
	db       UseCaseDB
}

func New(s Storage, ucs UseCaseShortURL) History {
	return History{storage: s, shortURL: ucs}
}

func (h *History) SetUseCaseDB(useCase UseCaseDB) {
	h.db = useCase
}

func (h History) hasUseCaseDB() bool {
	return h.db != nil
}

func (h History) Create(m *model.History) error {
	if !h.hasUseCaseDB() {
		return fmt.Errorf("the db dependency is required")
	}

	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()

	tx, err := h.db.Tx()
	if err != nil {
		return fmt.Errorf("h.db.Tx(): %v", err)
	}

	if err := h.storage.CreateWithTx(tx, m); err != nil {
		if errRollback := tx.Rollback(context.TODO()); errRollback != nil {
			return fmt.Errorf("h.storage.CreateWithTx(): rollback error %v, %w", errRollback, err)
		}

		return fmt.Errorf("h.storage.CreateWithTx(): %v", err)
	}

	if err := h.shortURL.IncrementTimes(tx, m.ShortURLID); err != nil {
		if errRollback := tx.Rollback(context.TODO()); errRollback != nil {
			return fmt.Errorf("h.shortURL.IncrementTimes(): rollback error %v, %w", errRollback, err)
		}

		return fmt.Errorf("h.shortURL.IncrementTimes(): %v", err)
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf("c.tx.Commit(): %v", err)
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
