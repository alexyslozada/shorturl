package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/user"
	"github.com/alexyslozada/shorturl/model"
)

type handler struct {
	useCase user.UseCase
	logger  *zap.SugaredLogger
}

func newHandler(uc user.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCase: uc, logger: l}
}

func (h handler) Create(c echo.Context) error {
	u := model.User{}
	err := c.Bind(&u)
	if err != nil {
		h.logger.Infow("can't bind user on create", "func", "Create", "internal", err)
		return c.JSON(http.StatusBadRequest, "Please verify the user structure")
	}

	err = h.useCase.Create(&u)
	if err != nil {
		h.logger.Errorw("can't create user", "func", "Create", "user", u, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't create the user")
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h handler) Delete(c echo.Context) error {
	ID := c.Param("id")
	uuidID, err := uuid.Parse(ID)
	if err != nil {
		h.logger.Infow("ID is not a valid uuid", "func", "Delete", "id", ID, "internal", err)
		return c.JSON(http.StatusBadRequest, "Please verify the ID is a valid uuid type")
	}

	err = h.useCase.Delete(uuidID)
	if err != nil {
		h.logger.Errorw("can't delete user", "func", "Delete", "id", uuidID, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't delete the user")
	}

	return c.JSON(http.StatusOK, nil)
}

func (h handler) ByEmail(c echo.Context) error {
	email := c.Param("email")
	u, err := h.useCase.ByEmail(email)
	if err != nil {
		h.logger.Errorw("can't get user by email", "func", "ByEmail", "email", email, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't get user")
	}

	return c.JSON(http.StatusOK, map[string]model.User{"data": u})
}

func (h handler) All(c echo.Context) error {
	us, err := h.useCase.All()
	if err != nil {
		h.logger.Errorw("can't get all users", "func", "All", "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't get users")
	}

	return c.JSON(http.StatusOK, map[string]model.Users{"data": us})
}
