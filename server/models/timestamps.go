package models

import (
	"time"

	"gorm.io/gorm"
)

// Timestamps provides CreatedAt, UpdatedAt, and soft-delete fields.
// Embed this in models that use a custom primary key instead of gorm.Model.
type Timestamps struct {
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
