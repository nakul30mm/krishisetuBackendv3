package models

import "time"

type Review struct {
	ID uint `gorm:"primaryKey"`

	EquipmentID *uint `json:"equipment_id" gorm:"index"`
	RentalID    *uint `json:"rental_id" gorm:"uniqueIndex"`

	ProductID            *uint `json:"product_id" gorm:"index"`
	MarketplaceRequestID *uint `json:"marketplace_request_id" gorm:"uniqueIndex"`

	ReviewerID uint `json:"reviewer_id" gorm:"not null"`

	Rating  int    `gorm:"not null"`
	Comment string `gorm:"type:text"`

	OwnerReply *string    `gorm:"type:text"`
	RepliedAt  *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
