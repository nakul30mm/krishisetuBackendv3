package models

import "gorm.io/gorm"

type QuestionVote struct {
	gorm.Model

	UserID     uint `gorm:"not null;index"`
	QuestionID uint `gorm:"not null;index"`

	Type string `gorm:"type:varchar(10);not null"` // "up" or "down"
}
