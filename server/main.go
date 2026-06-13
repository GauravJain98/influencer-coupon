package main

import (
	"flag"
	"gauravjain98/influencer-coupon/app"
	"log"
)

func main() {
	server := flag.Bool("server", false, "run the HTTP server")
	worker := flag.Bool("worker", false, "run the background worker")
	flag.Parse()

	if *server && *worker {
		log.Fatal("use only one of -server or -worker")
	}

	config := app.Config{}
	config.Load()

	if *worker {
		app.RunWorker(config)
		return
	}

	app.RunServer(config)
}
