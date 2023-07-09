package model

import (
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	Id        string `gorm:"primaryKey"`
	Sender    string `gorm:"type:varchar(36);not null"`
	Recipient string `gorm:"type:varchar(36);not null"`
	Message   string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
