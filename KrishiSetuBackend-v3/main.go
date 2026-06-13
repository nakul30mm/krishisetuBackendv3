package main

import (
	"krishisetu-backend/database"
	"krishisetu-backend/router"
	"krishisetu-backend/src/business"
	"krishisetu-backend/src/handlers"
	"krishisetu-backend/src/repository"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()
	db := database.DB

	// Repositories
	userRepo := repository.NewUserRepository(db)
	equipmentRepo := repository.NewEquipmentRepository(db)
	rentalRepo := repository.NewRentalRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	productRepo := repository.NewProductRepository(db)
	mrRepo := repository.NewMarketplaceRequestRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	answerRepo := repository.NewAnswerRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Services
	authService := business.NewAuthService(userRepo)
	equipmentService := business.NewEquipmentService(equipmentRepo, userRepo)
	rentalService := business.NewRentalService(rentalRepo, equipmentRepo)
	reviewService := business.NewReviewService(reviewRepo, rentalRepo, equipmentRepo, productRepo, mrRepo, userRepo)
	productService := business.NewProductService(productRepo, userRepo)
	mrService := business.NewMarketplaceRequestService(mrRepo, productRepo)
	qaService := business.NewQAService(questionRepo, answerRepo)
	adminService := business.NewAdminService(adminRepo, userRepo, productRepo, equipmentRepo, reviewRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService, rentalService.AutoCompleteRentals)
	rentalHandler := handlers.NewRentalHandler(rentalService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	productHandler := handlers.NewProductHandler(productService)
	mrHandler := handlers.NewMarketplaceRequestHandler(mrService)
	qaHandler := handlers.NewQAHandler(qaService)
	adminHandler := handlers.NewAdminHandler(adminService)
	uploadHandler := handlers.NewUploadHandler()

	// Router setup
	r := router.SetupRouter(
		authHandler,
		equipmentHandler,
		rentalHandler,
		reviewHandler,
		productHandler,
		mrHandler,
		qaHandler,
		adminHandler,
		uploadHandler,
	)

	// Run
	r.Run("0.0.0.0:8080")
}
