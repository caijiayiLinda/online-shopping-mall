package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/handlers"
	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	seed := flag.Bool("seed", false, "Seed the database with initial users")
	drop := flag.Bool("drop", false, "Drop all tables")
	flag.Parse()

	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Drop tables if flag is set
	if *drop {
		models.DropTables()
		return
	}

	// Seed users if flag is set
	if *seed {
		models.SeedUsers()
		return
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create database connection string with parseTime parameter
	dbConn := fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPassword, dbName)

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
  router.GET("/auth/csrf-token", authHandler.GetCSRFToken)
  router.GET("/auth/check", authHandler.CheckAuth)
  router.POST("/auth/login", authHandler.Login)
  router.POST("/auth/logout", authHandler.Logout)
  router.POST("/auth/register", authHandler.Register)
  router.POST("/auth/change-password", authHandler.AuthMiddleware(), authHandler.ChangePassword)
  
  // Protected routes
  adminGroup := router.Group("/admin")
  adminGroup.Use(authHandler.AdminAuthMiddleware())
  {
    adminGroup.POST("/products/create", productHandler.CreateProduct)
	adminGroup.POST("/products/update/:id", productHandler.UpdateProduct)
    adminGroup.DELETE("/products/delete/:id", productHandler.DeleteProduct)
    adminGroup.POST("/categories/create", categoryHandler.CreateCategory)
    adminGroup.POST("/categories/update", categoryHandler.UpdateCategory)
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
	router.POST("/checkout/paypal", gin.WrapH(handlers.CheckoutHandler(gormDB)))
	router.POST("/paypal/webhook", gin.WrapH(handlers.PayPalWebhookHandler(gormDB)))
	router.GET("/admin/orders", gin.WrapH(handlers.GetOrdersHandler(gormDB)))
	router.GET("/orders/by-email", gin.WrapH(handlers.GetRecentOrdersByEmailHandler(gormDB)))

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
