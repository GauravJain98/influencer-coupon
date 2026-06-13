package models

import "gorm.io/gorm"

type ChannelAffiliateVideo struct {
	gorm.Model
	Link               string `gorm:"type:text;primaryKey"`
	ChannelAffiliateID uint   `gorm:"primaryKey"`

	Video            Video            `gorm:"foreignKey:Link;references:Link"`
	ChannelAffiliate ChannelAffiliate `gorm:"foreignKey:ChannelAffiliateID"`
}
