package database

import (
  "context"
  "log"
  "os"
  "time"

  "github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {

  databaseUrl := os.Getenv("DATABASE_URL")

  config, err := pgxpool.ParseConfig(databaseUrl)
  if err != nil {
    log.Fatal(err)
  }

  config.MaxConns = 10
  config.MinConns = 2
  config.MaxConnLifetime = time.Hour

  pool, err := pgxpool.NewWithConfig(context.Background(), config)

  if err != nil {
    log.Fatal("Unable to connect to database:", err)
  }

  Pool = pool

  log.Println("Database connected")
}