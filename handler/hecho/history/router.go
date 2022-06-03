package history

import (
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"

	"github.com/alexyslozada/shorturl/handler/hecho/middleware"

	"github.com/alexyslozada/shorturl/domain/history"
)

const (
	path             = "/v1/histories"
	pathAll          = ""
	pathByShortURLID = "/short-url/:id"
)

func NewRouter(e *echo.Echo, uc history.UseCase, l *zap.SugaredLogger, middlewareFunc middleware.UseCase) {
	h := newHandler(uc, l)

	g := e.Group(path, middlewareFunc.ValidatePermission)
	g.GET(pathByShortURLID, h.ByShortURLID)
	g.GET(pathAll, h.All)
}
