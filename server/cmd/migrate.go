package main

import (
	"log"

	"github.com/GauravJain98/influencer-coupon/server/app"
	"github.com/GauravJain98/influencer-coupon/server/config"
	"github.com/GauravJain98/influencer-coupon/server/models"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	cfg := config.Config{}
	cfg.Load()

	db := app.DbConnect(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// for i := range channels {
	// 	utils.GetChannelIDAndNameFromHandle(ctx, cfg.YoutubeApiKey, &channels[i])
	// }

	// if err := seedChannels(db, channels); err != nil {
	// 	log.Fatal(err)
	// }

	db.AutoMigrate(
		&models.Channel{},
		&models.Affiliate{},
		&models.Video{},
		&models.ChannelAffiliate{},
		&models.ChannelAffiliateVideo{},
		&models.ScrapingError{},
	)
	// log.Printf("seeded %d channels", len(channels))
}
