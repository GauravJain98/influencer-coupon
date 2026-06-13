package models

import (
	"time"

	"gorm.io/gorm"
)

type Channel struct {
	ChannelID     string  `gorm:"type:text;column:channel_id"`
	Handle        *string `gorm:"type:text"`
	Name          *string `gorm:"type:text"`
	LastScrapedAt *time.Time
	gorm.Model

	Videos            []Video            `gorm:"foreignKey:ChannelID;references:ChannelID"`
	ChannelAffiliates []ChannelAffiliate `gorm:"foreignKey:ChannelID;references:ChannelID"`
}
