package models

import "gorm.io/gorm"

type MarketplaceRequest struct {
	gorm.Model

	ProductID uint     `json:"product_id"`
	Product   *Product `json:"product" gorm:"foreignKey:ProductID"`

	BuyerID uint  `json:"buyer_id"`
	Buyer   *User `json:"buyer" gorm:"foreignKey:BuyerID"`

	SellerID uint  `json:"seller_id"`
	Seller   *User `json:"seller" gorm:"foreignKey:SellerID"`

	RequestedPrice    float64 `json:"requested_price"`
	Quantity          int     `json:"quantity"`
	Status            string  `json:"status" gorm:"default:'PENDING'"`             // PENDING, APPROVED, REJECTED
	TransactionStatus string  `json:"transaction_status" gorm:"default:'PENDING'"` // YES, NO, PENDING
}
