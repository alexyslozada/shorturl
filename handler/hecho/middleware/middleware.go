package middleware

import "github.com/labstack/echo/v4"

type UseCase interface {
	ValidatePermission(next echo.HandlerFunc) echo.HandlerFunc
}
