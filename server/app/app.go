package app

import (
	"fmt"
	"log"

	"gauravjain98/influencer-coupon/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	//TODO: Handle allowed origins, cors, headers etc

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Streamer")
	})

	return router
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Channel{},
		&models.Affiliate{},
		&models.Video{},
		&models.ChannelAffiliate{},
		&models.ChannelAffiliateVideo{},
	)
}

func Run(config Config) {
	var db *gorm.DB
	var err error
	if config.DriverName == "sqlite3" {
		db, err = gorm.Open(sqlite.Open(config.SqlUrl), &gorm.Config{})
	} else if config.DriverName == "postgresql" {
		db, err = gorm.Open(postgres.Open(config.SqlUrl), &gorm.Config{})
	}

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := Migrate(db); err != nil {
		log.Fatal(err)
	}

	
	router := SetupRouter(db)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
