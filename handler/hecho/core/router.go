package core

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/domain/shorturl"
)

const (
	path = "/:short"
)

func NewRouter(e *echo.Echo, ucs shorturl.UseCase, uch history.UseCase, l *zap.SugaredLogger) {
	h := newHandler(ucs, uch, l)

	e.GET(path, h.Redirect)
}
