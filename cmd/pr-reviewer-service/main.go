package main

import (
	"context"
	"fmt"
	"log"
	//http2 "net/http"
	"os"

	app2 "github.com/blumgardt/pr-reviewer-service.git/internal/app"
	"github.com/blumgardt/pr-reviewer-service.git/internal/config"
	//"github.com/blumgardt/pr-reviewer-service.git/internal/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	configPath string = "config.toml"
)

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
	log.Printf("Connecting to DB: %s\n", connString)

	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal("Error ping db")
	}

	logger := log.New(os.Stdout, "[api]", log.Ldate|log.Ltime)

	app := app2.NewApp(cfg, logger, conn)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
