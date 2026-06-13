package models

import "gorm.io/gorm"

type Equipment struct {
	gorm.Model

	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	PricePerDay float64 `json:"price_per_day"`
	PriceUnit   string  `json:"price_unit"` // hour, day, week, etc.

	LocationPincode  string `json:"pincode"`
	LocationDistrict string `json:"district"`
	LocationState    string `json:"state"`
	LocationCity     string `json:"city"`
	LocationArea     string `json:"location"`

	Status     string  `json:"status"`
	BookedFrom *string `json:"booked_from"`
	BookedTo   *string `json:"booked_to"`

	OwnerID uint `json:"owner_id"`
	Owner   User `json:"owner" gorm:"foreignKey:OwnerID"`

	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
	Image3 string `json:"image3"`

	AverageRating float64 `gorm:"default:0" json:"average_rating"`
	TotalReviews  int     `gorm:"default:0" json:"total_reviews"`
}
