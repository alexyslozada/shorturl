package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/permission"
	"github.com/alexyslozada/shorturl/model"
)

type Middleware struct {
	useCasePermission permission.UseCase
	logger            *zap.SugaredLogger
	secretKey         string
}

func New(uc permission.UseCase, l *zap.SugaredLogger, sk string) Middleware {
	return Middleware{useCasePermission: uc, logger: l, secretKey: sk}
}

func (m *Middleware) SetPermission(uc permission.UseCase) {
	m.useCasePermission = uc
}

func (m *Middleware) hasPermission() bool {
	return m.useCasePermission != nil
}

func (m Middleware) ValidatePermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !m.hasPermission() {
			m.logger.Errorw(model.ErrMissingDependency.Error(), "func", "ValidatePermission", "internal", "has not use case")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": model.ErrMissingDependency.Error()})
		}

		token, err := m.authToken(c.Request().Header)
		if err != nil {
			m.logger.Errorw("couldn't get auth token", "func", "authToken", "internal", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		userID, err := uuid.Parse(token["user_id"].(string))
		if err != nil {
			m.logger.Errorw("couldn't parse uuid from token", "func", "uuid.Parse", "internal", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "el token no tiene el id del usuario"})
		}

		hasPermission, err := m.useCasePermission.HasPermission(userID, c.Request().Method)
		if err != nil {
			m.logger.Errorw("couldn't exec has permission", "func", "HasPermission", "internal", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if !hasPermission {
			m.logger.Infow("user has not permission", "userID", token["user_id"], "method", c.Request().Method)
			return c.JSON(http.StatusForbidden, map[string]string{"error": "you don't have access for this option"})
		}

		return next(c)
	}
}

func (m Middleware) authToken(header http.Header) (jwt.MapClaims, error) {
	tokenRequest := header.Get("Authorization")
	if len(tokenRequest) == 0 {
		return nil, fmt.Errorf("the auth token is empty")
	}

	if strings.Contains(tokenRequest, "Bearer ") {
		tokenRequest = tokenRequest[7:]
	}

	token, err := jwt.Parse(tokenRequest, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			m.logger.Errorw("Signing method is not valid", "func", "token.Method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// TODO implements a .pem cert file rather than model.Secret
		return []byte(m.secretKey), nil
	})
	if err != nil {
		m.logger.Errorw("couldn't parse token", "func", "jwt.Parse", "internal", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		m.logger.Infow(fmt.Sprintf("the type is %T", token.Claims))
		m.logger.Errorw("couldn't assertion parse claims", "func", "token.Claims", "parse", ok, "validToken", token.Valid)
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}
