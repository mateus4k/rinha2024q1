package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mateus4k/rinha2024q1/controller"
	"github.com/mateus4k/rinha2024q1/db"
)

func main() {
	ctx := context.Background()

	// Postgres
	config, err := pgxpool.ParseConfig(fmt.Sprintf("user=admin password=123 host=%s port=5432 dbname=rinha sslmode=disable pool_max_conns=%s", os.Getenv("DB_HOST"), os.Getenv("DB_POOL")))
	if err != nil {
		panic(err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	queries := db.New(conn)

	// API
	app := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	ctrl := controller.NewController(queries, conn)
	app.Get("/clientes/:id/extrato", ctrl.GetTransactions)
	app.Post("/clientes/:id/transacoes", ctrl.CreateTransaction)
	app.Listen(":3000")
}
