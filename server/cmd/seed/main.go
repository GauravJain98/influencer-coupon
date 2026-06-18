package main

import (
	"context"
	"log"
	"time"

	"github.com/GauravJain98/influencer-coupon/server/app"
	"github.com/GauravJain98/influencer-coupon/server/config"
	"github.com/GauravJain98/influencer-coupon/server/models"
	"github.com/GauravJain98/influencer-coupon/server/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cfg := config.Config{}
	cfg.Load()

	db := app.DbConnect(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	channels := []models.Channel{
		// Add channels here, for example:
		// {
		// 	Handle: stringPtr("@algorithm1313"),
		// },
		{
			Handle: stringPtr("@Moobub"),
		},
		// {
		// 	Handle: stringPtr("@WillTennyson"),
		// },
	}

	for i := range channels {
		err := utils.GetChannelIDAndNameFromHandle(ctx, cfg.YoutubeApiKey, &channels[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := seedChannels(db, channels); err != nil {
		log.Fatal(err)
	}

	app.Migrate(db)

	log.Printf("seeded %d channels", len(channels))
}

func seedChannels(db *gorm.DB, channels []models.Channel) error {
	if len(channels) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "channel_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"handle",
			"name",
			"last_scraped_at",
			"updated_at",
		}),
	}).Create(&channels).Error
}

func stringPtr(value string) *string {
	return &value
}
