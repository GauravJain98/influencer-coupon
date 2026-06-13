package main

import (
	"flag"
	"log"

	"github.com/GauravJain98/influencer-coupon/server/app"
	"github.com/GauravJain98/influencer-coupon/server/config"
)

func main() {
	server := flag.Bool("server", false, "run the HTTP server")
	worker := flag.Bool("worker", false, "run the background worker")
	flag.Parse()

	if *server && *worker {
		log.Fatal("use only one of -server or -worker")
	}

	cfg := config.Config{}
	cfg.Load()

	if *worker {
		app.RunWorker(cfg)
		return
	}

	app.RunServer(cfg)
}
