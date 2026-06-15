package models

import "gorm.io/gorm"

type ScrapingError struct {
	ErrorObject  map[string]interface{} `gorm:"serializer:json;type:text;column:error_object;not null"`
	FunctionName string                 `gorm:"type:text;column:function_name;not null;index"`
	Description  string                 `gorm:"type:text;column:description"`
	gorm.Model
}
