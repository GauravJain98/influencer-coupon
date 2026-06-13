package models

import "gorm.io/gorm"

type Video struct {
	Link            string  `gorm:"type:text;"`
	ChannelID       string  `gorm:"type:text;not null;index"`
	Title           *string `gorm:"type:text;column:title"`
	Description     *string `gorm:"type:text;column:description"`
	Status          int     `gorm:"type:integer;not null;default:1"`
	Evaluted        *bool   `gorm:"type:bool;column:evaluted"`
	NeedsRedo       bool    `gorm:"type:bool;not null;default:false;column:needs_redo"`
	EvaluationError *string `gorm:"type:text;column:evaluation_error"`
	gorm.Model

	Channel                Channel                 `gorm:"foreignKey:ChannelID;references:ChannelID"`
	ChannelAffiliateVideos []ChannelAffiliateVideo `gorm:"foreignKey:Link"`
}
