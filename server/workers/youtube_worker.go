package workers

import (
	"context"
	"log"
	"time"

	"github.com/GauravJain98/influencer-coupon/server/config"
	"github.com/GauravJain98/influencer-coupon/server/models"
	"github.com/GauravJain98/influencer-coupon/server/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewChannelBackfillWorker(config config.Config, db *gorm.DB) {
	errorFound := ""
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	channel, err := gorm.G[models.Channel](db).
		Where("backfill_completed_at IS NULL").
		Order("backfill_last_run_at ASC").
		First(ctx)
	channel.BackfillLastRunAt = time.Now().UTC()
	channel.LastScrapedAt = time.Now().UTC()

	log.Printf("ChannelID: %s ChannelName: %s ChannlBackfillLastRun: %s", channel.ChannelID, *channel.Name, channel.BackfillLastRunAt)
	if err != nil {
		// TODO: handle error
		log.Printf("Error: get channels needing backfill: %v", err)
		return
	}
	publishedBefore := time.Now()
	for true {
		videos, err := utils.GetChannelVideos(ctx, config.YoutubeApiKey, channel.ChannelID, publishedBefore)
		videoFoundCount := len(videos)
		if err != nil {
			errorFound = errorFound + "\n" + err.Error()
			log.Printf("Error: get channels needing backfill: %v", err)
			return
		}
		if videoFoundCount <= 0 {
			break
		}

		newVideos := make([]models.Video, 0, len(videos))
		for _, video := range videos {
			videoCount, err := gorm.G[models.Video](db).
				Where("link = ?", video.Link).Count(ctx, "link")
			if err != nil {
				errorFound = errorFound + "\n" + err.Error()
				log.Printf("Error checking video exists : %v", err)
				continue
			}

			if videoCount != 0 {
				// log.Printf("Video already exists, %d", videoCount)
				continue
			}
			newVideos = append(newVideos, video)
		}

		publishedBefore, newVideos, err = utils.ListVideoDetails(ctx, config.YoutubeApiKey, newVideos)
		if err != nil {
			errorFound = errorFound + "\n" + err.Error()
			log.Printf("Error getting video details : %v", err)
			continue
		}

		err = gorm.G[models.Video](db, clause.OnConflict{DoNothing: true}).CreateInBatches(ctx, &newVideos, 100)

		// err = gorm.G[models.Video](db).Create(ctx, &video)
		if err != nil {
			errorFound = errorFound + "\n" + err.Error()
			log.Printf("Error creating video details : %v", err)
			continue
		}
		log.Printf("Videos added")
	}

	// Save every found error in the db
	if errorFound != "" {
		channel.BackfillLastError = &errorFound
	}

	currentTime := time.Now().UTC()
	channel.BackfillCompletedAt = &currentTime

	_, err = gorm.G[models.Channel](db).
		Where("channel_id = ?", channel.ChannelID).
		Updates(ctx, channel)
	if err != nil {
		log.Printf("Error updating channel backfill run: %v", err)
		return
	}

}
