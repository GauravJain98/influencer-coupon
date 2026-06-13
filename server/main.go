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
	app.Run(config)
	fmt.Printf("HELLO WORLD")

}
