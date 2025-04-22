package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/models"
	"gorm.io/gorm"
)

type OrderResponse struct {
	ID           uint      `json:"id"`
	Currency     string    `json:"currency"`
	MerchantEmail string   `json:"merchant_email"`
	Salt         string    `json:"salt"`
	Products     []models.OrderProduct `json:"products"`
	TotalPrice   float64   `json:"total_price"`
	UserID       *uint     `json:"user_id"`
	Username     string    `json:"username"`
	Digest       string    `json:"digest"`
	Invoice      string    `json:"invoice"`
	Status       string    `json:"status"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

func GetOrdersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orders []models.Order
		
		if err := db.Preload("Products").Find(&orders).Error; err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		var response []OrderResponse
		for _, o := range orders {
			response = append(response, OrderResponse{
				ID:            o.ID,
				Currency:      o.Currency,
				MerchantEmail: o.MerchantEmail,
				Salt:          o.Salt,
				Products:      o.Products,
				TotalPrice:    o.TotalPrice,
				UserID:        o.UserID,
				Username:      o.Username,
				Digest:        o.Digest,
				Invoice:       o.Invoice,
				Status:        o.Status,
				CreatedAt:     o.CreatedAt.Format(time.RFC3339),
				UpdatedAt:     o.UpdatedAt.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetRecentOrdersByEmailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email parameter is required", http.StatusBadRequest)
			return
		}

		var orders []models.Order
		if err := db.Preload("Products").
			Where("merchant_email = ?", email).
			Order("created_at desc").
			Limit(5).
			Find(&orders).Error; err != nil {
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		// Convert to response format
		var response []OrderResponse
		for _, o := range orders {
			response = append(response, OrderResponse{
				ID:            o.ID,
				Currency:      o.Currency,
				MerchantEmail: o.MerchantEmail,
				Salt:          o.Salt,
				Products:      o.Products,
				TotalPrice:    o.TotalPrice,
				UserID:        o.UserID,
				Username:      o.Username,
				Digest:        o.Digest,
				Invoice:       o.Invoice,
				Status:        o.Status,
				CreatedAt:     o.CreatedAt.Format(time.RFC3339),
				UpdatedAt:     o.UpdatedAt.Format(time.RFC3339),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
