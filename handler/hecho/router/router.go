package router

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/domain/dbutil"
	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/domain/shorturl"
	"github.com/alexyslozada/shorturl/domain/user"
	routerCore "github.com/alexyslozada/shorturl/handler/hecho/core"
	routerHistory "github.com/alexyslozada/shorturl/handler/hecho/history"
	routerShortURL "github.com/alexyslozada/shorturl/handler/hecho/shorturl"
	routerUser "github.com/alexyslozada/shorturl/handler/hecho/user"
	"github.com/alexyslozada/shorturl/storage/postgres"
)

func Start(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger) {
	// H
	historyRouter(e, db, l)
	// R
	redirectRouter(e, db, l)
	// S
	shortURLRouter(e, db, l)
	// U
	userRouter(e, db, l)
}

func shortURLUseCase(db *pgxpool.Pool, l *zap.SugaredLogger) shorturl.ShortURL {
	storage := postgres.NewShortURL(db)
	useCase := shorturl.New(storage, l)
	useCase.SetUseCaseDB(dbutil.New(db))
	useCase.SetUseCaseHistory(history.New(postgres.NewHistory(db)))

	return useCase
}

func historyUseCase(db *pgxpool.Pool) history.UseCase {
	storage := postgres.NewHistory(db)
	useCase := history.New(storage)

	return useCase
}

func userRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger) {
	storage := postgres.NewUser(db)
	useCase := user.New(storage)
	routerUser.NewRouter(e, useCase, l)
}

func shortURLRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger) {
	useCase := shortURLUseCase(db, l)
	routerShortURL.NewRouter(e, useCase, l)
}

func historyRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger) {
	useCase := historyUseCase(db)
	routerHistory.NewRouter(e, useCase, l)
}

func redirectRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger) {
	shortURLUC := shortURLUseCase(db, l)
	historyUC := historyUseCase(db)
	routerCore.NewRouter(e, shortURLUC, historyUC, l)
}
