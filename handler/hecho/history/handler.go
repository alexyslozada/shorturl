package history

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/model"
)

type handler struct {
	useCase history.UseCase
	logger  *zap.SugaredLogger
}

func newHandler(uc history.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCase: uc, logger: l}
}

func (h handler) ByShortURLID(c echo.Context) error {
	ID := c.Param("id")
	uuidID, err := uuid.Parse(ID)
	if err != nil {
		h.logger.Infow("ID is not a valid uuid type", "func", "ByShortURLID", "id", ID, "internal", err)
		return c.JSON(http.StatusBadRequest, "ID has a not valid uuid format")
	}

	// If client wants filter by dates
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	if isFilteredByDates(from, to) {
		fromDate, toDate, err := h.parseFromAndToDates(from, to)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "From or To Date have a not valid format")
		}

		ss, err := h.useCase.ByShortURLIDAndDates(uuidID, fromDate, toDate)
		if err != nil {
			h.logger.Errorw("can't get history by short url ID and dates", "func", "ByShortURLIDAndDates", "id", uuidID, "internal", err)
			return c.JSON(http.StatusInternalServerError, "Ups!!! can't get history")
		}

		return c.JSON(http.StatusOK, map[string]model.Histories{"data": ss})
	}

	// If we fall here, means client want query only by ID
	ss, err := h.useCase.ByShortURLID(uuidID)
	if err != nil {
		h.logger.Errorw("can't get history by short url ID", "func", "ByShortURLID", "id", uuidID, "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't get history")
	}

	return c.JSON(http.StatusOK, map[string]model.Histories{"data": ss})
}

func (h handler) All(c echo.Context) error {
	ss, err := h.useCase.All()
	if err != nil {
		h.logger.Errorw("can't get all history", "func", "All", "internal", err)
		return c.JSON(http.StatusInternalServerError, "Ups!!! can't get all history")
	}

	return c.JSON(http.StatusOK, map[string]model.Histories{"data": ss})
}

func (h handler) parseFromAndToDates(from, to string) (time.Time, time.Time, error) {
	fromDate, err := time.Parse(time.RFC3339, from)
	if err != nil {
		h.logger.Infow("Date `From` has a not valid format", "func", "ByShortURLID", "from", from, "internal", err)
		return time.Time{}, time.Time{}, err
	}
	toDate, err := time.Parse(time.RFC3339, to)
	if err != nil {
		h.logger.Infow("Date `To` has a not valid format", "func", "ByShortURLID", "to", to, "internal", err)
		return time.Time{}, time.Time{}, err
	}

	return fromDate, toDate, nil
}

func isFilteredByDates(from, to string) bool {
	return from != "" && to != ""
}
