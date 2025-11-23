// @title           PR Reviewer Assignment Service API
// @version         1.0
// @description     Сервис назначения ревьюеров на Pull Request'ы.
// @BasePath        /
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	app2 "github.com/blumgardt/pr-reviewer-service.git/internal/app"
	"github.com/blumgardt/pr-reviewer-service.git/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/blumgardt/pr-reviewer-service.git/docs"
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

	conn, err := connectWithRetry(ctx, connString)
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

func connectWithRetry(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		log.Printf("Connecting to DB (attempt %d): %s", attempt, connString)

		pool, err = pgxpool.New(ctx, connString)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				log.Println("DB connection established")
				return pool, nil
			} else {
				err = pingErr
			}
		}

		log.Printf("DB not ready: %v", err)
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to DB after retries: %w", err)
}
