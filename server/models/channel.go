package models

import (
	"time"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	ChannelID           string             `gorm:"type:text;column:channel_id;unique;not null;check:length(channel_id) > 0"`
	Handle              *string            `gorm:"type:text"`
	Name                *string            `gorm:"type:text"`
	LastScrapedAt       time.Time          `gorm:"default:'2000-01-01 00:00:00'"` // The default date is 1st Jan 2000 so if no run has happened it is before the project
	BackfillLastRunAt   time.Time          `gorm:"column:backfill_last_run_at;not null;;default:'2000-01-01 00:00:00'"`
	BackfillCompletedAt *time.Time         `gorm:"column:backfill_completed_at;"`
	BackfillLastError   *string            `gorm:"type:text;column:backfill_last_error"`
	Videos              []Video            `gorm:"foreignKey:ChannelID;references:ChannelID"`
	ChannelAffiliates   []ChannelAffiliate `gorm:"foreignKey:ChannelID;references:ChannelID"`
}
