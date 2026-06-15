package workers

import (
	"context"
	"log"
	"time"

	"github.com/GauravJain98/influencer-coupon/server/config"
	"github.com/GauravJain98/influencer-coupon/server/models"
	"github.com/GauravJain98/influencer-coupon/server/utils"
	"gorm.io/gorm"
)

func NewChannelBackfillWorker(config config.Config, db *gorm.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// channel, err := gorm.G[models.Channel](db).
	// 	Where("backfill_completed_at IS NULL").First(ctx)

	channel, err := gorm.G[models.Channel](db).
		Where("backfill_completed_at IS NULL").
		Order("backfill_last_run_at ASC").
		First(ctx)

	channel.BackfillLastRunAt = time.Now()

	if err != nil {
		// TODO: handle error
		log.Printf("Error: get channels needing backfill: %v", err)
		return
	}
	videos, err := utils.GetChannelVideos(ctx, config.YoutubeApiKey, channel.ChannelID, nil)
	if err != nil {
		errMessage := err.Error()
		channel.BackfillLastError = &errMessage
		// TODO: improve the updateing
		_, updateErr := gorm.G[models.Channel](db).
			Where("channel_id = ?", channel.ChannelID).
			Updates(ctx, channel)
		if updateErr != nil {
			log.Printf("Error updating channel backfill error: %v", updateErr)
		}
		log.Printf("Error: get channels needing backfill: %v", err)
		return
	}
	




	channel.BackfillLastRunAt = time.Now()
	_, err = gorm.G[models.Channel](db).
		Where("channel_id = ?", channel.ChannelID).
		Updates(ctx, channel)
	if err != nil {
		log.Printf("Error updating channel backfill run: %v", err)
		return
	}

	log.Printf("ChannelID: %s ChannelName: %s ChannlBackfillLastRun", channel.ChannelID, *channel.Name, channel.BackfillLastRunAt)
}
