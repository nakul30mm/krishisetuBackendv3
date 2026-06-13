package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model

	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`

	Price    float64 `json:"price"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`

	LocationPincode  string `json:"pincode"`
	LocationDistrict string `json:"district"`
	LocationState    string `json:"state"`
	LocationCity     string `json:"city"`
	LocationArea     string `json:"location"`

	SellerID uint `json:"seller_id"`
	Seller   User `json:"seller" gorm:"foreignKey:SellerID"`

	AverageRating float64 `json:"average_rating" gorm:"default:0"`
	TotalReviews  int     `json:"total_reviews" gorm:"default:0"`

	Image1 string `json:"image1"`
	Image2 string `json:"image2"`
	Image3 string `json:"image3"`

	IsAvailable bool `gorm:"default:true"`
}
