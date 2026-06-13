package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenRouterToken string
	SecretKey       string
	SqlUrl          string
	DriverName      string
}

// type App struct {
// 	Config Config
// 	Db     *sql.DB
// }

func (config *Config) Load() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: unable to find a .env file")
	}

	config.SecretKey = os.Getenv("SECRET_KEY")

	if config.SecretKey == "" {
		log.Fatal("SECRET_KEY not set!")
	}

	config.SqlUrl = os.Getenv("SQL_URL")
	if config.SqlUrl == "" {
		log.Fatal("SQL_URL not set!")
	}

	config.DriverName = os.Getenv("DRIVER_NAME")
	if config.SqlUrl == "" {
		config.DriverName = "sqlite3"
	}

}
