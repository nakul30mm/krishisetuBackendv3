package models

type Location struct {
	Pincode  string `json:"pincode"`
	District string `json:"district"`
	State    string `json:"state"`
	City     string `json:"city"`
	Location string `json:"location"`
}
