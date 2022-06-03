package user

import (
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/handler/hecho/middleware"
	"github.com/labstack/echo/v4"

	"github.com/alexyslozada/shorturl/domain/user"
)

const (
	path        = "/v1/users"
	pathAll     = ""
	pathByID    = "/id/:id"
	pathByEmail = "/email/:email"
)

func NewRouter(e *echo.Echo, uc user.UseCase, l *zap.SugaredLogger, middlewareFunc middleware.UseCase) {
	h := newHandler(uc, l)

	g := e.Group(path, middlewareFunc.ValidatePermission)
	g.POST(pathAll, h.Create)
	g.DELETE(pathByID, h.Delete)
	g.GET(pathByEmail, h.ByEmail)
	g.GET(pathAll, h.All)
}
