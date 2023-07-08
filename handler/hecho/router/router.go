package router

import (
	"log"

	"go.uber.org/zap"

	"github.com/alexyslozada/shorturl/model"
	"github.com/alexyslozada/shorturl/storage/sheets"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/alexyslozada/shorturl/domain/history"
	"github.com/alexyslozada/shorturl/domain/login"
	"github.com/alexyslozada/shorturl/domain/permission"
	"github.com/alexyslozada/shorturl/domain/shorturl"
	"github.com/alexyslozada/shorturl/domain/user"
	routerCore "github.com/alexyslozada/shorturl/handler/hecho/core"
	routerHistory "github.com/alexyslozada/shorturl/handler/hecho/history"
	routerLogin "github.com/alexyslozada/shorturl/handler/hecho/login"
	"github.com/alexyslozada/shorturl/handler/hecho/middleware"
	routerShortURL "github.com/alexyslozada/shorturl/handler/hecho/shorturl"
	routerUser "github.com/alexyslozada/shorturl/handler/hecho/user"
	"github.com/alexyslozada/shorturl/storage/postgres"
)

func Start(e *echo.Echo, db *pgxpool.Pool, secretKey string, sheetsConfig model.Sheets, l *zap.SugaredLogger) {
	middlewarePermission := middlewareUseCase(db, l, secretKey)
	// H
	historyRouter(e, db, l, middlewarePermission)
	// L
	loginRouter(e, db, l, secretKey)
	// R
	redirectRouter(e, db, sheetsConfig, l)
	// S
	shortURLRouter(e, db, sheetsConfig, l, middlewarePermission)
	// U
	userRouter(e, db, l, middlewarePermission)
}

func shortURLUseCase(db *pgxpool.Pool, sheetsConfig model.Sheets, l *zap.SugaredLogger) shorturl.ShortURL {
	storage := postgres.NewShortURL(db)
	useCase := shorturl.New(storage, sheetsConfig, l)
	useCase.SetUseCaseHistory(history.New(postgres.NewHistory(db)))
	if sheetsConfig.HasToReportToSheets {
		sheetsService, err := sheets.New(sheetsConfig.ConfigFile, l)
		if err != nil {
			log.Fatalf("[ERROR]: cannot setup google sheets, error was: %v", err)
		}
		useCase.SetUseCaseSheets(sheetsService)
	}

	return useCase
}

func historyUseCase(db *pgxpool.Pool) history.UseCase {
	storage := postgres.NewHistory(db)
	useCase := history.New(storage)

	return useCase
}

func userRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger, middlewareFunc middleware.UseCase) {
	storage := postgres.NewUser(db)
	useCase := user.New(storage)
	routerUser.NewRouter(e, useCase, l, middlewareFunc)
}

func shortURLRouter(e *echo.Echo, db *pgxpool.Pool, sheetsConfig model.Sheets, l *zap.SugaredLogger, middlewareUseCase middleware.UseCase) {
	useCase := shortURLUseCase(db, sheetsConfig, l)
	routerShortURL.NewRouter(e, &useCase, l, middlewareUseCase)
}

func historyRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger, middlewareUseCase middleware.UseCase) {
	useCase := historyUseCase(db)
	routerHistory.NewRouter(e, useCase, l, middlewareUseCase)
}

func redirectRouter(e *echo.Echo, db *pgxpool.Pool, sheetsConfig model.Sheets, l *zap.SugaredLogger) {
	shortURLUC := shortURLUseCase(db, sheetsConfig, l)
	historyUC := historyUseCase(db)
	routerCore.NewRouter(e, &shortURLUC, historyUC, l)
}

func loginRouter(e *echo.Echo, db *pgxpool.Pool, l *zap.SugaredLogger, secretKey string) {
	storage := postgres.NewUser(db)
	useCase := login.New(storage, secretKey)
	routerLogin.NewRouter(e, useCase, l)
}

func middlewareUseCase(db *pgxpool.Pool, l *zap.SugaredLogger, secretKey string) middleware.UseCase {
	storage := postgres.NewPermission(db)
	useCase := permission.New(storage)
	return middleware.New(useCase, l, secretKey)
}
