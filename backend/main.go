package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create database connection string
	dbConn := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)

	// Initialize Gin router
	router := gin.Default()

	// Connect to MySQL database
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize GORM
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	productHandler := &handlers.ProductHandler{DB: db, Logger: log.Default()}
	categoryHandler := &handlers.CategoryHandler{DB: db, Logger: log.Default()}
	authHandler := &handlers.AuthHandler{DB: gormDB}

  // Setup routes
  router.POST("/login", authHandler.Login)
  router.POST("/logout", authHandler.Logout)
  
  // Protected routes
  adminGroup := router.Group("/admin")
  adminGroup.Use(authHandler.AdminAuthMiddleware())
  {
    adminGroup.POST("/products/create", productHandler.CreateProduct)
    adminGroup.PUT("/products/update", productHandler.UpdateProduct)
    adminGroup.DELETE("/products/delete", productHandler.DeleteProduct)
    adminGroup.POST("/categories/create", categoryHandler.CreateCategory)
    adminGroup.PUT("/categories/update", categoryHandler.UpdateCategory)
    adminGroup.DELETE("/categories/delete", categoryHandler.DeleteCategory)
  }

  // Public routes
  router.GET("/products", productHandler.ListProducts)
  router.GET("/products/:id", productHandler.GetProduct)
  router.GET("/products/category", productHandler.GetProductsByCategoryID)
  router.GET("/categories", categoryHandler.ListCategories)
  router.GET("/categories/id", categoryHandler.GetCategoryIDByName)
  router.GET("/categories/:id", categoryHandler.GetCategory)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Start server with graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}

// Middleware functions would go here...
