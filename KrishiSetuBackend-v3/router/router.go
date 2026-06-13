package router

import (
	"krishisetu-backend/src/handlers"
	"krishisetu-backend/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	equipmentHandler *handlers.EquipmentHandler,
	rentalHandler *handlers.RentalHandler,
	reviewHandler *handlers.ReviewHandler,
	productHandler *handlers.ProductHandler,
	marketplaceRequestHandler *handlers.MarketplaceRequestHandler,
	qaHandler *handlers.QAHandler,
	adminHandler *handlers.AdminHandler,
	uploadHandler *handlers.UploadHandler,
) *gin.Engine {
	r := gin.Default()

	// Static files
	r.Static("/uploads", "./uploads")

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth (public)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.POST("/forgot-password", authHandler.ForgotPassword)
	r.POST("/verify-otp", authHandler.VerifyOTP)
	r.POST("/reset-password", authHandler.ResetPassword)

	// Public rentals
	r.GET("/equipments", equipmentHandler.GetAllEquipments)
	r.GET("/equipments/:id/unavailable-dates", equipmentHandler.GetUnavailableDates)
	r.GET("/equipments/:id/reviews", reviewHandler.GetEquipmentReviews)

	// Public marketplace
	r.GET("/marketplace/products", productHandler.GetAllProducts)
	r.GET("/marketplace/products/:id/reviews", reviewHandler.GetProductReviews)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth(), middleware.BlockedCheck())
	{
		// Profile
		auth.GET("/profile", authHandler.Profile)
		auth.PUT("/profile", authHandler.UpdateProfile)

		// Equipments
		auth.POST("/equipments", equipmentHandler.CreateEquipment)
		auth.GET("/my/equipments", equipmentHandler.GetMyEquipments)
		auth.PUT("/equipments/:id", equipmentHandler.UpdateEquipment)
		auth.DELETE("/equipments/:id", equipmentHandler.DeleteEquipment)

		// Q&A — Questions
		auth.POST("/questions", qaHandler.CreateQuestion)
		auth.GET("/questions/community", qaHandler.GetCommunityQuestions)
		auth.GET("/questions/my", qaHandler.GetMyQuestions)
		auth.GET("/questions/replied", qaHandler.GetMyRepliedQuestions)

		// Q&A — Answers
		auth.POST("/answers", qaHandler.CreateAnswer)
		auth.GET("/questions/:id/answers", qaHandler.GetAnswersByQuestion)
		auth.GET("/answers/my", qaHandler.GetMyAnswers)

		// Q&A — Voting
		auth.POST("/questions/vote", qaHandler.VoteQuestion)
		auth.POST("/answers/vote", qaHandler.VoteAnswer)

		// Marketplace Requests
		auth.POST("/marketplace/requests", marketplaceRequestHandler.CreateMarketplaceRequest)
		auth.GET("/marketplace/requests/sent", marketplaceRequestHandler.GetSentMarketplaceRequests)
		auth.GET("/marketplace/requests/received", marketplaceRequestHandler.GetReceivedMarketplaceRequests)
		auth.POST("/marketplace/requests/:id/confirm-transaction", marketplaceRequestHandler.ConfirmMarketplaceTransaction)
		auth.POST("/marketplace/requests/:id/:action", marketplaceRequestHandler.UpdateMarketplaceRequestStatus)
		auth.DELETE("/marketplace/requests/:id", marketplaceRequestHandler.DeleteMarketplaceRequest)
		auth.PUT("/marketplace/requests/:id", marketplaceRequestHandler.UpdateMarketplaceRequest)

		// Rentals
		auth.POST("/rentals", rentalHandler.CreateRental)
		auth.GET("/rentals/sent", rentalHandler.GetMyRentalRequests)
		auth.GET("/rentals/received", rentalHandler.GetRentalRequestsForOwner)
		auth.POST("/rentals/:id/approve", rentalHandler.ApproveRental)
		auth.POST("/rentals/:id/reject", rentalHandler.RejectRental)
		auth.DELETE("/rentals/:id", rentalHandler.DeleteRental)
		auth.PUT("/rentals/:id", rentalHandler.UpdateRental)

		// Reviews
		auth.POST("/reviews", reviewHandler.CreateReview)
		auth.PUT("/reviews/:id", reviewHandler.UpdateReview)
		auth.DELETE("/reviews/:id", reviewHandler.DeleteReview)
		auth.PUT("/reviews/:id/reply", reviewHandler.ReplyToReview)
		auth.DELETE("/reviews/:id/reply", reviewHandler.DeleteReply)

		// Marketplace
		auth.POST("/marketplace/products", productHandler.CreateProduct)
		auth.GET("/marketplace/my/products", productHandler.GetMyProducts)
		auth.PUT("/marketplace/products/:id", productHandler.UpdateProduct)
		auth.DELETE("/marketplace/products/:id", productHandler.DeleteProduct)

		// Uploads
		auth.POST("/upload", uploadHandler.UploadImage)
	}

	// Admin-only routes
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.BlockedCheck(), middleware.AdminCheck())
	{
		admin.GET("/stats", adminHandler.GetAdminStats)
		admin.GET("/users", adminHandler.GetAdminUsers)
		admin.POST("/users/:id/block", adminHandler.BlockUser)
		admin.POST("/users/:id/unblock", adminHandler.UnblockUser)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)

		// Moderation
		admin.DELETE("/products/:id", adminHandler.AdminDeleteProduct)
		admin.DELETE("/equipments/:id", adminHandler.AdminDeleteEquipment)
		admin.DELETE("/reviews/:id", adminHandler.AdminDeleteReview)
	}

	return r
}
