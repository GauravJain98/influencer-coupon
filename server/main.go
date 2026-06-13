package main

import (
	"fmt"
	"gauravjain98/influencer-coupon/app"
)

func main() {
	config := app.Config{}
	config.Load()
	app.Run(config)
	fmt.Printf("HELLO WORLD")

}
