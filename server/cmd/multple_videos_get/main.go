package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GauravJain98/influencer-coupon/server/app"
	"github.com/GauravJain98/influencer-coupon/server/config"
	"github.com/GauravJain98/influencer-coupon/server/models"
	"github.com/GauravJain98/influencer-coupon/server/utils"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg := config.Config{}
	cfg.Load()

	db := app.DbConnect(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	channel, err := gorm.G[models.Channel](db).
		Where("backfill_completed_at IS NULL").
		Order("backfill_last_run_at ASC").
		First(ctx)

	videos, err := utils.GetChannelVideos(ctx, cfg.YoutubeApiKey, channel.ChannelID, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	_, videos, err = utils.ListVideoDetails(ctx, cfg.YoutubeApiKey, videos)
	for i, video := range videos {
		fmt.Println(*video.Title)
		fmt.Println(video.ChannelID)
		fmt.Println(video.Link)
		fmt.Println(*video.Description)
		if i == 3 {
			break
		}
	}

}
