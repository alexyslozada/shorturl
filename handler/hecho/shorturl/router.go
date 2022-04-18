package shorturl

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/shorturl"
	"github.com/alexyslozada/shorturl/handler/hecho/middleware"
)

const (
	path        = "/v1/short-urls"
	pathAll     = ""
	pathByID    = "/id/:id"
	pathByShort = "/short/:short"
)

func NewRouter(e *echo.Echo, uc shorturl.UseCase, l *zap.SugaredLogger, middlewareFunc middleware.UseCase) {
	h := newHandler(uc, l)

	g := e.Group(path, middlewareFunc.ValidatePermission)
	g.POST(pathAll, h.Create)
	g.PUT(pathAll, h.Update)
	g.DELETE(pathByID, h.Delete)
	g.GET(pathByShort, h.ByShort)
	g.GET(pathAll, h.All)
}
