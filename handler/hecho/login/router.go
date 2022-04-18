package login

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/login"
)

const (
	path    = "/v1/login"
	pathAll = ""
)

func NewRouter(e *echo.Echo, uc login.UseCase, l *zap.SugaredLogger) {
	h := newHandler(uc, l)

	g := e.Group(path)
	g.POST(pathAll, h.Login)
}
