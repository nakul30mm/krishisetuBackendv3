package models

import "gorm.io/gorm"

type AnswerVote struct {
	gorm.Model

	UserID   uint `gorm:"not null;index"`
	AnswerID uint `gorm:"not null;index"`

	Type string `gorm:"type:varchar(10);not null"` // "up" or "down"
}
