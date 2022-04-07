package core

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/dbutil"
	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/domain/shorturl"
	"github.com/alexyslozada/shorturl/model"
)

type handler struct {
	useCaseShortURL shorturl.UseCase
	useCaseHistory  history.UseCase
	useCaseDB       dbutil.UseCase
	logger          *zap.SugaredLogger
}

func newHandler(ucs shorturl.UseCase, uch history.UseCase, ucd dbutil.UseCase, l *zap.SugaredLogger) handler {
	return handler{useCaseShortURL: ucs, useCaseHistory: uch, useCaseDB: ucd, logger: l}
}

func (h handler) Redirect(c echo.Context) error {
	short := c.Param("short")
	shortURL, err := h.useCaseShortURL.ByShort(short)
	if err != nil {
		h.logger.Errorw("can't get short by short url", "func", "Redirect", "short", short, "internal", err)
		// We will return no content for this handler b/c this is used by a final client
		return c.JSON(http.StatusNoContent, "this url is not found")
	}

	go func() {
		tx, err := h.useCaseDB.Tx()
		if err != nil {
			h.logger.Errorw("couldn't get a transaction", "func", "Tx", "internal", err)
			return
		}
		// Is safe call rollback if the connection is not close
		defer func(tx pgx.Tx) {
			err := tx.Rollback(context.TODO())
			if err != nil {
				h.logger.Errorw("couldn't rollback transaction", "func", "Rollback", "internal", err)
			}
		}(tx)

		m := model.History{ShortURLID: shortURL.ID}
		err = h.useCaseHistory.CreateWithTx(tx, &m)
		if err != nil {
			h.logger.Errorw("couldn't create the history register", "func", "CreateWithTx", "short", short, "internal", err)
			return
		}

		err = tx.Commit(context.TODO())
		if err != nil {
			h.logger.Errorw("couldn't commit transaction", "func", "CreateWithTx", "short", short, "internal", err)
		}
	}()

	return c.Redirect(http.StatusMovedPermanently, shortURL.RedirectTo)
}
