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

	// Initialize handlers with logger
	productHandler := &handlers.ProductHandler{DB: db, Logger: log.Default()}
	categoryHandler := &handlers.CategoryHandler{DB: db, Logger: log.Default()}

	// Setup routes
	http.HandleFunc("/products", productHandler.ListProducts)
	http.HandleFunc("/products/create", productHandler.CreateProduct) 
	http.HandleFunc("/products/update", productHandler.UpdateProduct)
	http.HandleFunc("/products/delete", productHandler.DeleteProduct)
	http.HandleFunc("/products/", productHandler.GetProduct)
	http.HandleFunc("/categories", categoryHandler.ListCategories)
	http.HandleFunc("/categories/create", categoryHandler.CreateCategory)
	http.HandleFunc("/categories/update", categoryHandler.UpdateCategory)
	http.HandleFunc("/categories/delete", categoryHandler.DeleteCategory)
	http.HandleFunc("/categories/", categoryHandler.GetCategory)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain middleware and handle root route
	http.Handle("/", logRequest(cors(http.DefaultServeMux)))

	// Start server with graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
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

// logRequest middleware logs incoming requests
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for multipart requests to avoid consuming request body
		if r.Header.Get("Content-Type") == "multipart/form-data" {
			log.Printf("Request: %s %s [multipart]", r.Method, r.URL.Path)
		} else {
			log.Printf("Request: %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

// cors middleware adds CORS headers
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
