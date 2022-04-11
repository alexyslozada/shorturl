package shorturl

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

type ShortURL struct {
	storage        Storage
	logger         Logger
	useCaseHistory UseCaseHistory
	useCaseDB      UseCaseDB
}

func New(s Storage, logger Logger) ShortURL {
	return ShortURL{storage: s, logger: logger}
}

func (s *ShortURL) SetUseCaseHistory(useCase UseCaseHistory) {
	s.useCaseHistory = useCase
}

func (s ShortURL) hasUseCaseHistory() bool {
	return s.useCaseHistory != nil
}

func (s *ShortURL) SetUseCaseDB(useCase UseCaseDB) {
	s.useCaseDB = useCase
}

func (s ShortURL) hasUseCaseDB() bool {
	return s.useCaseDB != nil
}

func (s ShortURL) Create(m *model.ShortURLRequest) error {
	if !strings.Contains(m.RedirectTo, HTTPProtocol) {
		return model.ErrWrongRedirect
	}

	short := model.ShortURL{
		ID:          uuid.New(),
		Short:       m.Short,
		RedirectTo:  m.RedirectTo,
		Description: m.Description,
		CreatedAt:   time.Now().Unix(),
	}

	if m.IsRandom {
		short.Short = randomPATH()
		m.Short = short.Short
	}

	return s.storage.Create(&short)
}

func (s ShortURL) Update(m *model.ShortURL) error {
	m.UpdatedAt = time.Now().Unix()

	return s.storage.Update(m)
}

func (s ShortURL) IncrementTimes(tx pgx.Tx, ID uuid.UUID) error {
	return s.storage.IncrementTimes(tx, ID)
}

func (s ShortURL) Delete(ID uuid.UUID) error {
	return s.storage.Delete(ID)
}

func (s ShortURL) ByShort(shortURL string) (model.ShortURL, error) {
	m, err := s.storage.ByShort(shortURL)
	if err != nil {
		return model.ShortURL{}, fmt.Errorf("s.storage.ByShort(): %w", err)
	}

	return m, nil
}

func (s ShortURL) ByShortToRedirect(shortURL string) (model.ShortURL, error) {
	m, err := s.storage.ByShort(shortURL)
	if err != nil {
		return model.ShortURL{}, fmt.Errorf("s.storage.ByShort(): %w", err)
	}

	go func() {
		if err := s.createHistoryAndIncrementTimes(m); err != nil {
			s.logger.Errorw("s.createHistoryAndIncrementTimes(): %v", err)
		}
	}()

	return m, nil
}

func (s ShortURL) All() (model.ShortURLs, error) {
	return s.storage.All()
}

func (s ShortURL) createHistoryAndIncrementTimes(shortURL model.ShortURL) error {
	if err := s.validateDependencies(); err != nil {
		return err
	}

	tx, err := s.useCaseDB.Tx()
	if err != nil {
		return fmt.Errorf("s.useCaseDB.Tx(): %w", err)
	}

	defer func(tx pgx.Tx) {
		if errRollback := tx.Rollback(context.TODO()); errRollback != nil {
			if errors.Is(errRollback, pgx.ErrTxClosed) {
				return
			}

			s.logger.Errorw("could not be rollback on tx", "internal", fmt.Errorf("tx.Rollback(): rollback error %v, %w", errRollback, err))
		}
	}(tx)

	m := model.History{ShortURLID: shortURL.ID}
	if err := s.useCaseHistory.CreateWithTx(tx, &m); err != nil {
		return fmt.Errorf("h.storage.CreateWithTx(): %v", err)
	}

	if err := s.IncrementTimes(tx, m.ShortURLID); err != nil {
		return fmt.Errorf("h.shortURL.IncrementTimes(): %v", err)
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return fmt.Errorf("c.tx.Commit(): %v", err)
	}

	return nil
}

func (s ShortURL) validateDependencies() error {
	if !s.hasUseCaseHistory() {
		return fmt.Errorf("%w: %s", model.ErrMissingDependency, "history")
	}

	if !s.hasUseCaseDB() {
		return fmt.Errorf("%w: %s", model.ErrMissingDependency, "db")
	}

	return nil
}

func randomPATH() string {
	resp := make([]rune, MaxLetters)
	lenAllowedLetters := len(allowedLetters)

	for i := range resp {
		resp[i] = allowedLetters[rand.Intn(lenAllowedLetters)]
	}

	return string(resp)
}
