package models

type ChannelAffiliateVideo struct {
	Link               string `gorm:"type:text;primaryKey"`
	ChannelAffiliateID uint   `gorm:"primaryKey"`
	Timestamps

	Video            Video            `gorm:"foreignKey:Link;references:Link"`
	ChannelAffiliate ChannelAffiliate `gorm:"foreignKey:ChannelAffiliateID"`
}
