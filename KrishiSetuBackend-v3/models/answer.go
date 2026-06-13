package models

import "gorm.io/gorm"

type Answer struct {
	gorm.Model

	Content string `json:"content" gorm:"type:text;not null"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`

	QuestionID uint     `json:"question_id"`
	Question   Question `json:"question"`

	Upvotes   int64 `json:"upvotes" gorm:"default:0"`
	Downvotes int64 `json:"downvotes" gorm:"default:0"`
}
