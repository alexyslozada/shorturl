// This package is a main file for Echo router and Postgres db
package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

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

	router.Start(e, dbPool, logger.Sugar())

	err := e.Start(":" + os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Println("Error: Couldn't start the server", err)
		os.Exit(1)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
		os.Exit(1)
	}
}

func getEcho() *echo.Echo {
	e := echo.New()
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
		log.Println("Can't connect to postgres db", err)
		os.Exit(1)
	}

	return db
}

func getLog() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println("Error: Couldn't create the logger", err)
		os.Exit(1)
	}

	return logger
}

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]time.Time{"data": time.Now()})
}
