package main

import (
	"fmt"
	"gauravjain98/influencer-coupon/app"
)

// "database/sql"
// "fmt"
// "log"

// "github.com/gin-gonic/gin"
// _ "github.com/mattn/go-sqlite3"

func main() {
	config := app.Config{}
	config.Load()
	
	fmt.Printf("HELLO WORLD")
	// db, err := sql.Open("sqlite3", "./foo.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// router := gin.Default()
	// // router.GET("/albums", getAlbums)

	// router.Run("localhost:8080")
}
