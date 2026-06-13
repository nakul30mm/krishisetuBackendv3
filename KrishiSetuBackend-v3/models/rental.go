package models

import "gorm.io/gorm"

type Rental struct {
	gorm.Model

	EquipmentID uint      `json:"equipment_id"`
	Equipment   Equipment `json:"equipment" gorm:"foreignKey:EquipmentID"`

	RenterID uint `json:"renter_id"`
	Renter   User `json:"renter" gorm:"foreignKey:RenterID"`

	OwnerID uint `json:"owner_id"`

	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`

	Status string `json:"status"`
}
