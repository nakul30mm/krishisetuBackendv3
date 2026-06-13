package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model

	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user"`

	Upvotes   int64 `json:"upvotes"`
	Downvotes int64 `json:"downvotes"`

	RepliesCount int64 `json:"replies_count" gorm:"column:replies_count"`
}
