package db

import (
	"log"

	"github.com/chimas/GoProject/config"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DBConnection() (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		// log.Fatal("Error loading .env file")
	}
	db, err := sqlx.ConnectContext(ctx, "postgres", config.LoadEnv().DB_URL)

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return nil, err
	}

	return db, nil
}
