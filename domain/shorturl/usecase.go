package shorturl

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/alexyslozada/shorturl/model"
)

type ShortURL struct {
	storage Storage
}

func New(s Storage) ShortURL {
	return ShortURL{storage: s}
}

func (s ShortURL) Create(m *model.ShortURLRequest) error {
	if !strings.Contains(m.RedirectTo, HTTPProtocol) {
		return ErrWrongRedirect
	}

	short := model.ShortURL{
		ID:          uuid.New(),
		Short:       m.Short,
		RedirectTo:  m.RedirectTo,
		Description: m.Description,
		CreatedAt:   time.Now().Unix(),
	}

	if m.IsRandom {
		m.Short = randomPATH()
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
	return s.storage.ByShort(shortURL)
}

func (s ShortURL) All() (model.ShortURLs, error) {
	return s.storage.All()
}

func randomPATH() string {
	resp := make([]rune, MaxLetters)
	lenAllowedLetters := len(allowedLetters)

	for i := range resp {
		resp[i] = allowedLetters[rand.Intn(lenAllowedLetters)]
	}

	return string(resp)
}
