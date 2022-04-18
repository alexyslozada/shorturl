package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/login"
	"github.com/alexyslozada/shorturl/model"
)

type handler struct {
	useCase login.UseCase
	logger  *zap.SugaredLogger
}

func newHandler(uc login.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCase: uc, logger: l}
}

func (h handler) Login(c echo.Context) error {
	m := model.LoginRequest{}
	err := c.Bind(&m)
	if err != nil {
		h.logger.Infow("the login request struct is not valid", "func", "Login", "internal", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "revisa la estructura enviada"})
	}

	token, err := h.useCase.Login(m.Email, m.Password)
	if err != nil {
		h.logger.Infow("error on login", "func", "Login", "internal", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "usuario o contraseña no válidos"})
	}

	return c.JSON(http.StatusOK, map[string]string{"data": token})
}
