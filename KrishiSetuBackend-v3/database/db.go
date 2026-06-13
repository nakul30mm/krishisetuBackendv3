package database

import (
	"fmt"
	"log"
	"os"

	"krishisetu-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if username == "" {
		username = "root"
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if dbname == "" {
		dbname = "krishisetu_db"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("MySQL connected successfully")

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Equipment{},
		&models.Question{},
		&models.Answer{},
		&models.QuestionVote{},
		&models.AnswerVote{},
		&models.MarketplaceRequest{}, // new Phase 1.5 feature
		&models.Rental{},
		&models.Review{},
		&models.Product{},
	); err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}
}
