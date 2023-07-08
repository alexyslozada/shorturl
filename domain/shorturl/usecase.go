package shorturl

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/model"
)

type ShortURL struct {
	storage        Storage
	sheetsConfig   model.Sheets
	logger         *zap.SugaredLogger
	useCaseHistory UseCaseHistory
	useCaseSheets  UseCaseSheets
}

func New(s Storage, sheetsConfig model.Sheets, logger *zap.SugaredLogger) ShortURL {
	return ShortURL{storage: s, sheetsConfig: sheetsConfig, logger: logger}
}

func (s *ShortURL) SetUseCaseHistory(useCase UseCaseHistory) {
	s.useCaseHistory = useCase
}

func (s *ShortURL) SetUseCaseSheets(useCase UseCaseSheets) {
	s.useCaseSheets = useCase
}

func (s *ShortURL) hasUseCaseHistory() bool {
	return s.useCaseHistory != nil
}

func (s *ShortURL) hasUseCaseSheets() bool {
	return s.useCaseSheets != nil
}

func (s *ShortURL) Create(m *model.ShortURL) error {
	if !strings.Contains(m.RedirectTo, HTTPProtocol) {
		return model.ErrWrongRedirect
	}

	m.ID = uuid.New()
	m.CreatedAt = time.Now().Unix()

	if m.IsRandom {
		m.Short = randomPATH()
	}

	return s.storage.Create(m)
}

func (s *ShortURL) Update(m *model.ShortURL) error {
	m.UpdatedAt = time.Now().Unix()

	return s.storage.Update(m)
}

func (s *ShortURL) Delete(ID uuid.UUID) error {
	return s.storage.Delete(ID)
}

func (s *ShortURL) ByShort(shortURL string) (model.ShortURL, error) {
	m, err := s.storage.ByShort(shortURL)
	if err != nil {
		return model.ShortURL{}, fmt.Errorf("s.storage.ByShort(): %w", err)
	}

	return m, nil
}

func (s *ShortURL) ByShortToRedirect(shortURL string) (model.ShortURL, error) {
	m, err := s.storage.ByShort(shortURL)
	if err != nil {
		return model.ShortURL{}, fmt.Errorf("s.storage.ByShort(): %w", err)
	}

	go s.createHistoryAndIncrementTimes(m)

	return m, nil
}

func (s *ShortURL) All(limit, offset int) (model.ShortURLs, error) {
	return s.storage.All(limit, offset)
}

func (s *ShortURL) createHistoryAndIncrementTimes(shortURL model.ShortURL) {
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

	if s.sheetsConfig.HasToReportToSheets {
		if !s.hasUseCaseSheets() {
			s.logger.Errorw("we need a use case sheets in order to report")
			return
		}

		err := s.useCaseSheets.AddRow(&shortURL, m.CreatedAt, s.sheetsConfig.SpreadsheetID)
		if err != nil {
			s.logger.Errorw(fmt.Sprintf("s.useCaseSheets.AddRow(shortURL, %d, %s", m.CreatedAt, s.sheetsConfig.SpreadsheetID))
		}
	}
}

func (s *ShortURL) validateDependencies() error {
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
