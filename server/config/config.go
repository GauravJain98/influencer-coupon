package config

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
	YoutubeApiKey   string
}

type envVar struct {
	Name     string
	Target   *string
	Required bool
	Default  string
}

func (config *Config) Load() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: unable to find a .env file")
	}

	vars := []envVar{
		{Name: "SECRET_KEY", Target: &config.SecretKey, Required: true},
		{Name: "SQL_URL", Target: &config.SqlUrl, Required: true},
		{Name: "DRIVER_NAME", Target: &config.DriverName, Default: "sqlite3"},
		{Name: "YOUTUBE_API_KEY", Target: &config.YoutubeApiKey, Required: true},
	}

	for _, env := range vars {
		value := os.Getenv(env.Name)

		if value == "" {
			value = env.Default
		}

		if value == "" && env.Required {
			log.Fatalf("%s not set!", env.Name)
		}

		*env.Target = value
	}

}
