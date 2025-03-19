package handlers

import (
	"bytes"
	"database/sql"
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
	"github.com/gin-gonic/gin"
	"backend/models"
)

type ProductHandler struct {
	DB     *sql.DB
	Logger *log.Logger
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Log request details
	h.Logger.Printf("CreateProduct request received")
	h.Logger.Printf("Request Headers: %v", c.Request.Header)

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		h.Logger.Printf("Error parsing form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}

	// Get form values
	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	name := form.Value["name"][0]
	price, err := strconv.ParseFloat(form.Value["price"][0], 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	description := form.Value["description"][0]

	// Handle file upload
	h.Logger.Printf("Attempting to get uploaded file from form-data")
	fileHeader, err := c.FormFile("image")
	if err != nil {
		h.Logger.Printf("Error getting uploaded file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to get file: %v", err)})
		return
	}
	file, err := fileHeader.Open()
	defer file.Close()
	h.Logger.Printf("Successfully received file: %s (%d bytes)", fileHeader.Filename, fileHeader.Size)

	// Create unique filename
	ext := filepath.Ext(fileHeader.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Ensure images directory exists
	if err := os.MkdirAll("/home/caijiayi/online-shopping-mall/public/images", 0755); err != nil {
		h.Logger.Printf("Error creating images directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create images directory"})
		return
	}

	// Save original file
	filePath := filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "original_"+newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		h.Logger.Printf("Error creating file %s: %v", filePath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}
	defer dst.Close()

	// Create a TeeReader to read the file once and write to both destination and memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Printf("Error reading file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
		return
	}

	if _, err := dst.Write(fileBytes); err != nil {
		h.Logger.Printf("Error writing file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// Verify file type
	contentType := http.DetectContentType(fileBytes)
	h.Logger.Printf("Detected file content type: %s", contentType)
	if !strings.Contains(contentType, "image/") {
		h.Logger.Printf("Invalid file type: %s", contentType)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid file type: %s. Only image files are allowed.", contentType)})
		return
	}

	// Decode image from memory
	h.Logger.Printf("Attempting to decode image, size: %d bytes", len(fileBytes))
	img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		h.Logger.Printf("Failed to decode image: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "The image file appears to be corrupted or in an unsupported format. Please try with a valid JPG, PNG or GIF image."})
		return
	}
	h.Logger.Printf("Successfully decoded image, dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var thumbnailPath string
	if width <= 300 && height <= 300 {
		// Use original image as thumbnail if it's small enough
		thumbnailPath = filePath
	} else {
		// Create thumbnail for larger images
	thumbnail := imaging.Resize(img, 300, 300, imaging.Lanczos)
	thumbnailFilename := "thumbnail_" + newFilename
	thumbnailPath = filepath.Join("/home/caijiayi/online-shopping-mall/public/images", thumbnailFilename)
	err = imaging.Save(thumbnail, thumbnailPath)
	if err != nil {
		h.Logger.Printf("Error saving thumbnail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save thumbnail"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	h.Logger.Print(result)

	// Get inserted product ID
	productID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
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
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		h.Logger.Printf("Error parsing form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}

	// Get form values
	categoryID, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	name := form.Value["name"][0]
	price, err := strconv.ParseFloat(form.Value["price"][0], 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}

	description := form.Value["description"][0]

	// Handle file upload if exists
	var newFilename string
	fileHeader, err := c.FormFile("image")
	if err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to open file: %v", err)})
			return
		}
		defer file.Close()

		// Create unique filename
		ext := filepath.Ext(fileHeader.Filename)
		newFilename = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		// Save original file
		filePath := filepath.Join("/home/caijiayi/online-shopping-mall/public/images", "original_"+newFilename)
		dst, err := os.Create(filePath)
		if err != nil {
			h.Logger.Printf("Error creating file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}
		defer dst.Close()

		// Read file into memory
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			h.Logger.Printf("Error reading file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
			return
		}

		if _, err := dst.Write(fileBytes); err != nil {
			h.Logger.Printf("Error writing file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}

		// Decode image from memory
		img, err := imaging.Decode(bytes.NewReader(fileBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process image"})
			return
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		var thumbnailPath string
		if width <= 300 && height <= 300 {
			// Use original image as thumbnail if it's small enough
			thumbnailPath = filePath
		} else {
			// Create thumbnail for larger images
			thumbnail := imaging.Resize(img, 300, 300, imaging.Lanczos)
			thumbnailFilename := "thumbnail_" + newFilename
			thumbnailPath = filepath.Join("/home/caijiayi/online-shopping-mall/public/images", thumbnailFilename)
			err = imaging.Save(thumbnail, thumbnailPath)
			if err != nil {
				h.Logger.Printf("Error saving thumbnail: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save thumbnail"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Return updated product with full image URLs
	product := models.Product{
		ID:          productID,
		CategoryID:  categoryID,
		Name:        name,
		Price:       price,
		Description: description,
		ImageURL:    "/images/" + newFilename,
		ThumbnailURL: "/images/thumbnail_" + newFilename,
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Delete product from database
	query := `DELETE FROM products WHERE pid = ?`
	result, err := h.DB.Exec(query, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Query product from database
	row := h.DB.QueryRow("SELECT * FROM products WHERE pid = ?", productID)
	var p models.Product
	err = row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) GetProductsByCategoryID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Query("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	rows, err := h.DB.Query("SELECT * FROM products WHERE catid = ?", categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	h.Logger.Printf("Handling ListProducts request")
	
	// Query products from database
	rows, err := h.DB.Query("SELECT pid, catid, name, price, description, image_url, thumbnail_url FROM products")
	if err != nil {
		h.Logger.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.Price, &p.Description, &p.ImageURL, &p.ThumbnailURL)
		if err != nil {
			h.Logger.Printf("Database scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		h.Logger.Printf("Database rows error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, products)
}
