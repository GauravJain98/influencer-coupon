package models

import "gorm.io/gorm"

type ChannelAffiliate struct {
	gorm.Model
	Link        *string `gorm:"type:text"`
	Code        *string `gorm:"type:text"`
	AffiliateID string  `gorm:"type:varchar(10);not null;index"`
	ChannelID   string  `gorm:"type:text;not null;index"`

	Affiliate              Affiliate               `gorm:"foreignKey:AffiliateID"`
	Channel                Channel                 `gorm:"foreignKey:ChannelID;references:ChannelID"`
	ChannelAffiliateVideos []ChannelAffiliateVideo `gorm:"foreignKey:ChannelAffiliateID"`
}
