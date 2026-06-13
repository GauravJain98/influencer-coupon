package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	//TODO: Handle allowed origins, cors, headers etc

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Streamer")
	})

	return router
}

func Run(config Config) {
	db, err := sql.Open(config.DriverName, config.SqlUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := SetupRouter(db)
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}

}
