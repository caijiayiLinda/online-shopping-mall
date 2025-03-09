package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend/models"
)

type ProductHandler struct {
	DB     *sql.DB
	Logger *log.Logger
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Log request details
	h.Logger.Printf("CreateProduct request received")
	h.Logger.Printf("Request Headers: %v", r.Header)
	h.Logger.Printf("Request Form Data: %v", r.Form)

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		h.Logger.Printf("Error parsing form: %v", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	// Handle file upload
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create unique filename
	ext := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Ensure images directory exists
	if err := os.MkdirAll("/home/caijiayi/online-shopping-mall/public/images", 0755); err != nil {
		http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
		return
	}

	// Save file
	filePath := filepath.Join("/home/caijiayi/online-shopping-mall/public/images", newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Log file save success
	h.Logger.Printf("Successfully saved image to: %s", filePath)

	// Insert product into database
		query := `INSERT INTO products (catid, name, price, description, image_url) VALUES (?, ?, ?, ?, ?)`
		result, err := h.DB.Exec(query, categoryID, name, price, description, "/images/" + newFilename)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get inserted product ID
	productID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Return created product with full image URL path
	product := models.Product{
		ID:          int(productID),
		CategoryID:  categoryID,
		Name:        name,
		Price:       price,
		Description: description,
		ImageURL:    "/images/" + newFilename,
	}

	// Log successful response
	h.Logger.Printf("Successfully created product: %+v", product)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		h.Logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Error creating response", http.StatusInternalServerError)
	}
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get product ID from URL
	productID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get form values
	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	// Handle file upload if exists
	var newFilename string
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Create unique filename
		ext := filepath.Ext(handler.Filename)
		newFilename = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		// Save file
		dst, err := os.Create(filepath.Join("/home/caijiayi/online-shopping-mall/public/images", newFilename))
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
	}

	// Update product in database
	var query string
	var result sql.Result
	if newFilename != "" {
		query = `UPDATE products SET catid = ?, name = ?, price = ?, description = ?, image_url = ? WHERE pid = ?`
		result, err = h.DB.Exec(query, categoryID, name, price, description, "/images/" + newFilename, productID)
	} else {
		query = `UPDATE products SET catid = ?, name = ?, price = ?, description = ? WHERE pid = ?`
		result, err = h.DB.Exec(query, categoryID, name, price, description, productID)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return updated product with full image URL path
	product := models.Product{
		ID:          productID,
		CategoryID:  categoryID,
		Name:        name,
		Price:       price,
		Description: description,
		ImageURL:    "/images/" + newFilename,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	productID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete product from database
	query := `DELETE FROM products WHERE pid = ?`
	result, err := h.DB.Exec(query, productID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	productID, err := strconv.Atoi(r.URL.Path[len("/products/"):])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Query product from database
	row := h.DB.QueryRow("SELECT pid, catid, name, price, description, image_url FROM products WHERE pid = ?", productID)

	var p models.Product
	err = row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) GetProductsByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(r.URL.Query().Get("category_id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	rows, err := h.DB.Query("SELECT pid, catid, name, price, description, image_url FROM products WHERE catid = ?", categoryID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	h.Logger.Printf("Handling ListProducts request")
	// Query products from database
	rows, err := h.DB.Query("SELECT pid, catid, name, price, description, image_url FROM products")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
