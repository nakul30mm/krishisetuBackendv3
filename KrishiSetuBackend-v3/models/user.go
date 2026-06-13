package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	FullName string `json:"full_name"`
	Email    string `json:"email" gorm:"unique"`
	Phone    string `json:"phone"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`

	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`

	Pincode  string `json:"pincode"`
	Location string `json:"location"`
	Password string `json:"-"`

	ProfilePicture string `json:"profile_picture"`

	AverageRating float64 `gorm:"default:0" json:"average_rating"`
	TotalRatings  int     `gorm:"default:0" json:"total_ratings"`

	ResetOTP         *string    `json:"-"`
	ResetOTPExpiry   *time.Time `json:"-"`
	ResetOTPVerified bool       `json:"-"`
	IsAdmin          bool       `gorm:"default:false" json:"is_admin"`
	IsBlocked        bool       `gorm:"default:false" json:"is_blocked"`
}
