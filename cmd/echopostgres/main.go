// This package is a main file for Echo router and Postgres db
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexyslozada/shorturl/model"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/alexyslozada/shorturl/handler/hecho/router"
	"github.com/alexyslozada/shorturl/storage/postgres"
)

func main() {
	loadEnv()

	logger := getLog()
	defer func() {
		_ = logger.Sync()
	}()

	e := getEcho()
	e.GET("/health", health)

	dbPool := getPostgres()
	defer dbPool.Close()

	err := postgres.Migrate(dbPool)
	if err != nil {
		log.Fatalln("Error: couldn't migrate database", err)
	}

	// TODO change this to a cert .pem
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) == 0 {
		panic("The JWT_SECRET_KEY env var must be set")
	}

	sheetsConf := getSheetsConfig()
	router.Start(e, dbPool, secretKey, sheetsConf, logger.Sugar())

	err = e.Start(":" + os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatalln("Error: Couldn't start the server", err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}
}

func getEcho() *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	return e
}

func getPostgres() *pgxpool.Pool {
	min := 3
	max := 100

	minConns := os.Getenv("DB_MIN_CONNS")
	maxConns := os.Getenv("DB_MAX_CONNS")

	if minConns != "" {
		v, err := strconv.Atoi(minConns)
		if err != nil {
			log.Println("Warning: DB_MIN_CONNS has not a valid value, we will set min connections to", min)
		} else {
			min = v
		}
	}
	if maxConns != "" {
		v, err := strconv.Atoi(maxConns)
		if err != nil {
			log.Println("Warning: DB_MAX_CONNS has not a valid value, we will set max connections to", max)
		} else {
			max = v
		}
	}

	db, err := postgres.New(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
		int32(min),
		int32(max),
	)
	if err != nil {
		log.Fatalln("Can't connect to postgres db", err)
	}

	return db
}

func getLog() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("Error: Couldn't create the logger", err)
	}

	return logger
}

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]time.Time{"data": time.Now()})
}

func getSheetsConfig() model.Sheets {
	willStore := os.Getenv("SHEETS_HAS_TO_STORE")
	if willStore == "" {
		return model.Sheets{}
	}

	hasToStore, err := strconv.ParseBool(willStore)
	if err != nil {
		log.Printf("[WARN] the %s is not a valid boolean value, will continue without store in google sheets", willStore)
		return model.Sheets{}
	}

	m := model.Sheets{
		HasToReportToSheets: hasToStore,
	}

	m.ConfigFile = os.Getenv("SHEETS_CONFIGURATION")
	if m.ConfigFile == "" {
		log.Fatalln("[ERROR] config file is mandatory")
	}
	m.SpreadsheetID = os.Getenv("SHEETS_SPREADSHEET_ID")
	if m.SpreadsheetID == "" {
		log.Fatalln("[ERROR] spreadsheet ID is mandatory")
	}
	m.SpreadsheetSheet = os.Getenv("SHEETS_SPREADSHEET_SHEET")
	if m.SpreadsheetSheet == "" {
		log.Fatalln("[ERROR] spreadsheet sheet is mandatory")
	}

	return m
}
