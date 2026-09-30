package main

import (
	"context"
	"log"
	"os"

	"github.com/Fulim13/author-api/internal/handler"
	"github.com/Fulim13/author-api/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer pool.Close()

	authorService := service.NewAuthorService(pool)
	authorHandler := handler.NewAuthorHandler(authorService)

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	handler.Routes(e, authorHandler)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
