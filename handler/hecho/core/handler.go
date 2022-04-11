package core

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/domain/shorturl"
)

type handler struct {
	useCaseShortURL shorturl.UseCase
	useCaseHistory  history.UseCase
	logger          *zap.SugaredLogger
}

func newHandler(ucs shorturl.UseCase, uch history.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCaseShortURL: ucs, useCaseHistory: uch, logger: l}
}

func (h handler) Redirect(c echo.Context) error {
	short := c.Param("short")
	shortURL, err := h.useCaseShortURL.ByShortToRedirect(short)
	if errors.Is(err, pgx.ErrNoRows) {
		return c.JSON(http.StatusNoContent, "this url is not found")
	}
	if err != nil {
		h.logger.Errorw("can't get short by short url", "func", "Redirect", "short", short, "internal", err)
		// We will return no content for this handler b/c this is used by a final client
		return c.JSON(http.StatusInternalServerError, "can't get short by short url")
	}

	return c.Redirect(http.StatusMovedPermanently, shortURL.RedirectTo)
}
