package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	REDIS_URL string
	DB_URL    string
}

func LoadEnv() EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		// log.Fatal("Error loading .env file")
	}

	redis_url := os.Getenv("REDIS_URL")
	db_url := os.Getenv("DB_URL")
	// redis_db := os.Getenv("REDIS_DB")
	// parsed_redis_db, err := strconv.Atoi(redis_db)
	// if err != nil {
	// 	panic("cannot parse redis DB number")
	// }

	return EnvVars{
		REDIS_URL: redis_url,
		DB_URL:    db_url,
	}
}
