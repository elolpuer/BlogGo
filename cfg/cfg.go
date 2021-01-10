package cfg

import (
	"log"
	"os"

	"fmt"

	"github.com/joho/godotenv"
)

func load(){
	var err error
	err = godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env", err.Error())
	}
}

func GetPostgres() string {
	load()
	PgHost := os.Getenv("PgHost")
	PgPort := os.Getenv("PgPort")
	PgUser := os.Getenv("PgUser")
	PgPass := os.Getenv("PgPass")
	PgDB := os.Getenv("PgDB")
	SSLmode := os.Getenv("SSLmode")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s ", PgHost, PgPort, PgUser, PgPass, PgDB, SSLmode)
}

func GetSessionKey() string {
	load()
	SessionKey := os.Getenv("SESSION_KEY")
	return SessionKey
}