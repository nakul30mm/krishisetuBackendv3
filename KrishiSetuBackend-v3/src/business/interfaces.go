package business

import "krishisetu-backend/models"

// UserRepository defines methods for user data access.
type UserRepository interface {
    // Create adds a new user record.
    Create(user *models.User) error
    // FindByEmail retrieves a user by email.
    FindByEmail(email string) (*models.User, error)
    // FindByID retrieves a user by ID.
    FindByID(id uint) (*models.User, error)
    // Update saves changes to a user.
    Update(user *models.User) error
    // UpdateFields updates specific fields of a user.
    UpdateFields(id uint, updates map[string]interface{}) error
    GetAll() ([]models.User, error)
    Delete(id uint) error
    FindByIdentifier(identifier string) (*models.User, error)
}

// AdminRepository defines methods for admin-specific data aggregation.
type AdminRepository interface {
    GetStats() (userCount int64, productCount int64, equipmentCount int64, requestCount int64, err error)
}

// EquipmentRepository defines methods for equipment data access.
type EquipmentRepository interface {
    // Create adds a new equipment record.
    Create(equipment *models.Equipment) error
    // GetFiltered retrieves equipment records matching filter options.
    GetFiltered(search, sort, minPrice, maxPrice, userPincode string) ([]models.Equipment, error)
    // GetByOwnerID retrieves all equipment owned by a user.
    GetByOwnerID(ownerID uint) ([]models.Equipment, error)
    // FindByID retrieves equipment by its ID (without preloads by default, or with Owner).
    FindByID(id uint) (*models.Equipment, error)
    // FindByIDWithOwner retrieves equipment by its ID, preloading the Owner.
    FindByIDWithOwner(id uint) (*models.Equipment, error)
    // Update saves changes to equipment.
    Update(equipment *models.Equipment) error
    // Delete removes an equipment record.
    Delete(id uint) error
    // CountActiveRentals returns the count of approved rentals for an equipment.
    CountActiveRentals(equipmentID uint) (int64, error)
    // GetApprovedRentals retrieves all approved rentals for an equipment.
    GetApprovedRentals(equipmentID uint) ([]models.Rental, error)
    // GetLatestApprovedRental retrieves the latest approved rental for an equipment.
    GetLatestApprovedRental(equipmentID uint) (*models.Rental, error)
}

// RentalRepository defines methods for rental data access.
type RentalRepository interface {
    Create(rental *models.Rental) error
    FindByID(id uint) (*models.Rental, error)
    FindByIDWithEquipmentAndOwner(id uint) (*models.Rental, error)
    Update(rental *models.Rental) error
    Delete(id uint) error
    GetByRenterID(renterID uint) ([]models.Rental, error)
    GetByOwnerID(ownerID uint) ([]models.Rental, error)
    GetApprovedRentalsForEquipment(equipmentID uint) ([]models.Rental, error)
    GetApprovedRentalsForEquipmentExcluding(equipmentID uint, excludingID uint) ([]models.Rental, error)
    CountApprovedRentalsForEquipment(equipmentID uint) (int64, error)
    RejectOverlappingPendingRentals(equipmentID uint, ignoreRentalID uint, startDate, endDate string) error
    GetAllApprovedRentals() ([]models.Rental, error)
    UpdateEquipmentStatus(equipmentID uint, status string) error
    GetReviewForRental(rentalID uint) (*models.Review, error)
}

// ReviewRepository defines methods for review data access and aggregation.
type ReviewRepository interface {
    Create(review *models.Review) error
    FindByID(id uint) (*models.Review, error)
    Update(review *models.Review) error
    Delete(id uint) error
    GetByEquipmentID(equipmentID uint) ([]models.Review, error)
    GetByProductID(productID uint) ([]models.Review, error)
    GetByRentalID(rentalID uint) (*models.Review, error)
    GetByMarketplaceRequestID(marketplaceRequestID uint) (*models.Review, error)
    UpdateEquipmentRating(equipmentID uint) error
    UpdateProductRating(productID uint) error
    UpdateOwnerRating(ownerID uint) error
}

// ProductRepository defines methods for product data access.
type ProductRepository interface {
    Create(product *models.Product) error
    GetFiltered(search, sort, minPrice, maxPrice, userPincode string) ([]models.Product, error)
    GetBySellerID(sellerID uint) ([]models.Product, error)
    FindByID(id uint) (*models.Product, error)
    FindByIDWithSeller(id uint) (*models.Product, error)
    Update(product *models.Product) error
    Delete(id uint) error
}

// MarketplaceRequestRepository defines methods for marketplace request data access.
type MarketplaceRequestRepository interface {
    Create(req *models.MarketplaceRequest) error
    FindByID(id uint) (*models.MarketplaceRequest, error)
    FindByIDWithProductAndSeller(id uint) (*models.MarketplaceRequest, error)
    Update(req *models.MarketplaceRequest) error
    Delete(id uint) error
    GetByBuyerID(buyerID uint, search string) ([]models.MarketplaceRequest, error)
    GetBySellerID(sellerID uint, search string) ([]models.MarketplaceRequest, error)
}

// QuestionRepository defines methods for Q&A questions data access.
type QuestionRepository interface {
    Create(question *models.Question) error
    FindByID(id uint) (*models.Question, error)
    FindByIDWithUser(id uint) (*models.Question, error)
    Update(question *models.Question) error
    GetCommunityQuestions(userID uint, search string) ([]models.Question, error)
    GetMyQuestions(userID uint, search string) ([]models.Question, error)
    GetMyRepliedQuestions(userID uint, search string) ([]models.Question, error)
    GetVote(userID uint, questionID uint) (*models.QuestionVote, error)
    CreateVote(vote *models.QuestionVote) error
    DeleteVote(vote *models.QuestionVote) error
    UpdateVote(vote *models.QuestionVote) error
    UpdateVotesCount(questionID uint, upvotes, downvotes int64) error
}

// AnswerRepository defines methods for Q&A answers data access.
type AnswerRepository interface {
    Create(answer *models.Answer) error
    FindByID(id uint) (*models.Answer, error)
    Update(answer *models.Answer) error
    GetByQuestionID(questionID uint) ([]models.Answer, error)
    GetByUserID(userID uint) ([]models.Answer, error)
    GetVote(userID uint, answerID uint) (*models.AnswerVote, error)
    CreateVote(vote *models.AnswerVote) error
    DeleteVote(vote *models.AnswerVote) error
    UpdateVote(vote *models.AnswerVote) error
    UpdateVotesCount(answerID uint, upvotes, downvotes int64) error
}
