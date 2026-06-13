package dto

// Admin DTOs
type AdminStatsDTO struct {
	TotalUsers      int64 `json:"total_users"`
	TotalProducts   int64 `json:"total_products"`
	TotalEquipments int64 `json:"total_equipments"`
	TotalRequests   int64 `json:"total_requests"`
}

// Auth DTOs
type RegisterDTO struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
	Pincode  string `json:"pincode"`
	Location string `json:"location"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileDTO struct {
	FullName       string `json:"full_name"`
	Phone          string `json:"phone"`
	Age            string `json:"age"`
	Gender         string `json:"gender"`
	City           string `json:"city"`
	District       string `json:"district"`
	State          string `json:"state"`
	Pincode        string `json:"pincode"`
	Location       string `json:"location"`
	ProfilePicture string `json:"profile_picture"`
}

// Equipment DTOs
type LocationDTO struct {
	Pincode  string `json:"pincode"`
	District string `json:"district"`
	State    string `json:"state"`
	City     string `json:"city"`
	Location string `json:"location"`
}

type CreateEquipmentDTO struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
	PricePerDay float64     `json:"price_per_day"`
	PriceUnit   string      `json:"price_unit"`
	Location    LocationDTO `json:"location"`
	Image1      string      `json:"image1"`
	Image2      string      `json:"image2"`
	Image3      string      `json:"image3"`
}

type UpdateEquipmentDTO struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	PricePerDay float64     `json:"price_per_day"`
	Location    LocationDTO `json:"location"`
	PriceUnit   string      `json:"price_unit"`
	Image1      string      `json:"image1"`
	Image2      string      `json:"image2"`
	Image3      string      `json:"image3"`
}

// Rental DTOs
type CreateRentalDTO struct {
	EquipmentID uint   `json:"equipment_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type UpdateRentalDTO struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Review DTOs
type CreateReviewDTO struct {
	EquipmentID     *uint  `json:"equipment_id"`
	RentalID        *uint  `json:"rental_id"`
	ProductID       *uint  `json:"product_id"`
	MarketRequestID *uint  `json:"market_request_id"`
	Rating          int    `json:"rating" binding:"required"`
	Comment         string `json:"comment"`
}

type UpdateReviewDTO struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

type OwnerReplyDTO struct {
	Reply string `json:"reply" binding:"required"`
}

// Product DTOs
type CreateProductDTO struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Unit        string      `json:"unit"`
	Quantity    float64     `json:"quantity"`
	Location    LocationDTO `json:"location"`
	Image1      string      `json:"image1"`
	Image2      string      `json:"image2"`
	Image3      string      `json:"image3"`
}

type UpdateProductDTO struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Unit        string      `json:"unit"`
	Quantity    float64     `json:"quantity"`
	Location    LocationDTO `json:"location"`
	Image1      string      `json:"image1"`
	Image2      string      `json:"image2"`
	Image3      string      `json:"image3"`
}

// Marketplace Request DTOs
type CreateMarketplaceRequestDTO struct {
	ProductID      uint    `json:"product_id"`
	RequestedPrice float64 `json:"requested_price"`
	Quantity       int     `json:"quantity"`
}

type UpdateMarketplaceRequestDTO struct {
	Quantity int `json:"quantity"`
}

type ConfirmTransactionDTO struct {
	Status string `json:"status" binding:"required"`
}

// Q&A DTOs
type CreateQuestionDTO struct {
	Content string `json:"content"`
}

type CreateAnswerDTO struct {
	Content    string `json:"content"`
	QuestionID uint   `json:"question_id"`
}

type VoteQuestionDTO struct {
	QuestionID uint   `json:"question_id"`
	Type       string `json:"type"`
}

type VoteAnswerDTO struct {
	AnswerID uint   `json:"answer_id"`
	Type     string `json:"type"`
}
