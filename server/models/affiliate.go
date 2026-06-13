package models

import "gorm.io/gorm"

type Affiliate struct {
	ID          string  `gorm:"type:varchar(10);primaryKey"`
	Name        *string `gorm:"type:text"`
	Description *string `gorm:"type:text"`
	gorm.Model

	ChannelAffiliates []ChannelAffiliate `gorm:"foreignKey:AffiliateID"`
}
