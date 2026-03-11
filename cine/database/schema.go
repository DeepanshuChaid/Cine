package database

import (
  "context"
  "log"
)

func InitSchema() {

  query := `
  CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL
  );
  `

  _, err := Pool.Exec(context.Background(), query)

  if err != nil {
    log.Fatal("Schema creation failed:", err)
  }

  log.Println("Schema ready")
}