package history

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/history"
)

const (
	path             = "/v1/histories"
	pathAll          = ""
	pathByShortURLID = "/short-url/:id"
)

func NewRouter(e *echo.Echo, uc history.UseCase, l *zap.SugaredLogger) {
	h := newHandler(uc, l)

	g := e.Group(path)
	g.GET(pathByShortURLID, h.ByShortURLID)
	g.GET(pathAll, h.All)
}
