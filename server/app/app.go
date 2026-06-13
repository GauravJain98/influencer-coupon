package app

import (
	"fmt"
	"log"

	"gauravjain98/influencer-coupon/models"
	"gauravjain98/influencer-coupon/routes"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	//TODO: Handle allowed origins, cors, headers etc

	routes.SetupAdminRoutes(router, db)
	routes.SetupUserRoutes(router, db)
	routes.SetupPublicRoutes(router, db)

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

func DbConnect(config Config) *gorm.DB {
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

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("Can not ping db", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatal(err)
	}
	return db
}

func RunServer(config Config) {

	db := DbConnect(config)

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	router := SetupRouter(db)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}

func RunWorker(config Config) {
	db := DbConnect(config)

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	log.Fatal("THIS HAS NOT BEEN IMPLEMENTED YET")
}
