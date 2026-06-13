package models

import "time"

type Channel struct {
	ChannelID     string  `gorm:"type:text;primaryKey;column:channel_id"`
	Handle        *string `gorm:"type:text"`
	Name          *string `gorm:"type:text"`
	LastScrapedAt *time.Time
	Timestamps

	Videos            []Video            `gorm:"foreignKey:ChannelID;references:ChannelID"`
	ChannelAffiliates []ChannelAffiliate `gorm:"foreignKey:ChannelID;references:ChannelID"`
}
