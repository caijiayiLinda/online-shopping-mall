package models

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DropTables() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Migrator().DropTable(&User{})
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	log.Println("Successfully dropped all tables")
}

func SeedUsers() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create database connection string
	dsn := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed admin user
	adminPassword := "Admin@1234"
	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := User{
		Email:    "admin@example.com",
		Password: string(adminHash),
		Admin:    true,
	}

	result := db.Create(&admin)
	if result.Error != nil {
		log.Fatalf("Failed to create admin user: %v", result.Error)
	}

	// Seed regular user
	userPassword := "User@1234"
	userHash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash user password: %v", err)
	}

	user := User{
		Email:    "user@example.com",
		Password: string(userHash),
		Admin:    false,
	}

	result = db.Create(&user)
	if result.Error != nil {
		log.Fatalf("Failed to create regular user: %v", result.Error)
	}

	log.Println("Successfully seeded admin and regular users")
	log.Println("Admin credentials: admin@example.com / Admin@1234")
	log.Println("User credentials: user@example.com / User@1234")
}
