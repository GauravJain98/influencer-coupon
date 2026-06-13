package models

type Affiliate struct {
	ID          string  `gorm:"type:varchar(10);primaryKey"`
	Name        *string `gorm:"type:text"`
	Description *string `gorm:"type:text"`
	Timestamps

	ChannelAffiliates []ChannelAffiliate `gorm:"foreignKey:AffiliateID"`
}
