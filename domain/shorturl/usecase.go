package shorturl

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alexyslozada/shorturl/model"
)

type ShortURL struct {
	storage        Storage
	logger         Logger
	useCaseHistory UseCaseHistory
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

	go s.createHistoryAndIncrementTimes(m)

	return m, nil
}

func (s ShortURL) All() (model.ShortURLs, error) {
	return s.storage.All()
}

func (s ShortURL) createHistoryAndIncrementTimes(shortURL model.ShortURL) {
	if err := s.validateDependencies(); err != nil {
		s.logger.Errorw(fmt.Sprintf("s.validateDependencies(): %v", err))
		return
	}

	m := model.History{ShortURLID: shortURL.ID}
	if err := s.useCaseHistory.Create(&m); err != nil {
		s.logger.Errorw(fmt.Sprintf("s.useCaseHistory.Create(): %v", err))
	}

	if err := s.storage.IncrementTimes(m.ShortURLID); err != nil {
		s.logger.Errorw(fmt.Sprintf("s.storage.IncrementTimes(%d): %v", m.ShortURLID, err))
	}
}

func (s ShortURL) validateDependencies() error {
	if !s.hasUseCaseHistory() {
		return fmt.Errorf("%w: %s", model.ErrMissingDependency, "history")
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
