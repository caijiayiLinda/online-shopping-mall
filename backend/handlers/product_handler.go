package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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
	h.Logger.Printf("Attempting to get uploaded file from form-data")
	file, handler, err := r.FormFile("image")
	if err != nil {
		h.Logger.Printf("Error getting uploaded file: %v", err)
		http.Error(w, fmt.Sprintf("Unable to get file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()
	h.Logger.Printf("Successfully received file: %s (%d bytes)", handler.Filename, handler.Size)

	// Create unique filename
	ext := filepath.Ext(handler.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Ensure images directory exists
	if err := os.MkdirAll("/home/caijiayi/online-shopping-mall/public/images", 0755); err != nil {
		h.Logger.Printf("Error creating images directory: %v", err)
		http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
		return
	}

	// Save original file
	filePath := filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "original_"+newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		h.Logger.Printf("Error creating file %s: %v", filePath, err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Create a TeeReader to read the file once and write to both destination and memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Printf("Error reading file: %v", err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	if _, err := dst.Write(fileBytes); err != nil {
		h.Logger.Printf("Error writing file: %v", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Verify file type
	contentType := http.DetectContentType(fileBytes)
	h.Logger.Printf("Detected file content type: %s", contentType)
	if !strings.Contains(contentType, "image/") {
		h.Logger.Printf("Invalid file type: %s", contentType)
		http.Error(w, fmt.Sprintf("Invalid file type: %s. Only image files are allowed.", contentType), http.StatusBadRequest)
		return
	}

	// Decode image from memory
	h.Logger.Printf("Attempting to decode image, size: %d bytes", len(fileBytes))
	img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		h.Logger.Printf("Failed to decode image: %v", err)
		http.Error(w, "The image file appears to be corrupted or in an unsupported format. Please try with a valid JPG, PNG or GIF image.", http.StatusBadRequest)
		return
	}
	h.Logger.Printf("Successfully decoded image, dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var thumbnailPath string
	if width <= 500 && height <= 500 {
		// Use original image as thumbnail if it's small enough
		thumbnailPath = filePath
	} else {
		// Create thumbnail for larger images
	thumbnail := imaging.Resize(img, 500, 500, imaging.Lanczos)
	thumbnailFilename := "thumbnail_" + newFilename
	thumbnailPath = filepath.Join("/home/caijiayi/online-shopping-mall/public/images", thumbnailFilename)
	err = imaging.Save(thumbnail, thumbnailPath)
	if err != nil {
		h.Logger.Printf("Error saving thumbnail: %v", err)
		http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
		return
	}
	h.Logger.Printf("Successfully created thumbnail: %s", thumbnailFilename)
	}

	// Log file save success
	h.Logger.Printf("Successfully saved image to: %s", filePath)
	h.Logger.Printf("Successfully saved thumbnail to: %s", thumbnailPath)

	// Insert product into database
	imageURL := "/images/" + "original_" + newFilename
	thumbnailURL := "/images/" + "thumbnail_" + newFilename
	h.Logger.Printf("Inserting product with image_url: %s, thumbnail_url: %s", imageURL, thumbnailURL)
	query := `INSERT INTO products (catid, name, price, description, image_url, thumbnail_url) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := h.DB.Exec(query, categoryID, name, price, description, imageURL, thumbnailURL)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	h.Logger.Print(result)

	// Get inserted product ID
	productID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Return created product with full image URLs
	product := models.Product{
		ID:          int(productID),
		CategoryID:  categoryID,
		Name:        name,
		Price:       price,
		Description: description,
		ImageURL:    "/images/original_" + newFilename,
		ThumbnailURL: "/images/thumbnail_" + newFilename,
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
	productID, err := strconv.Atoi(r.URL.Query().Get("pid"))
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

		// Save original file
		dst, err := os.Create(filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "original_"+newFilename))
		if err != nil {
			h.Logger.Printf("Error creating file: %v", err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Read file into memory
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			h.Logger.Printf("Error reading file: %v", err)
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}

		if _, err := dst.Write(fileBytes); err != nil {
			h.Logger.Printf("Error writing file: %v", err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		// Decode image from memory
		img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, "Unable to process image", http.StatusInternalServerError)
		return
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var thumbnailPath string
	if width <= 500 && height <= 500 {
		// Use original image as thumbnail if it's small enough
		thumbnailPath = filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "thumbnail_"+newFilename)
	} else {
		// Create thumbnail for larger images
		thumbnail := imaging.Resize(img, 500, 500, imaging.Lanczos)
		thumbnailPath = filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "thumbnail_"+newFilename)
		err = imaging.Save(thumbnail, thumbnailPath)
		if err != nil {
			http.Error(w, "Unable to save thumbnail", http.StatusInternalServerError)
			return
		}
	}
	}

	// Update product in database
	var query string
	var result sql.Result
	if newFilename != "" {
		query = `UPDATE products SET catid = ?, name = ?, price = ?, description = ?, image_url = ?, thumbnail_url = ? WHERE pid = ?`
		result, err = h.DB.Exec(query, categoryID, name, price, description, "/images/" + newFilename, "/images/thumbnail_" + newFilename, productID)
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
		ThumbnailURL: "/images/thumbnail_" + newFilename,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	productID, err := strconv.Atoi(r.URL.Query().Get("pid"))
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
	row := h.DB.QueryRow("SELECT * FROM products WHERE pid = ?", productID)
	var p models.Product
	err = row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
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

	rows, err := h.DB.Query("SELECT * FROM products WHERE catid = ?", categoryID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
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
	rows, err := h.DB.Query("SELECT pid, catid, name, price, description, image_url, thumbnail_url FROM products")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
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
