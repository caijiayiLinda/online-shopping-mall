package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/plutov/paypal/v4"
	"gorm.io/gorm"
	"backend/models"
)

type PayPalOrderRequest struct {
	CartItems []CartItem `json:"cartItems" binding:"required,min=1"`
	Invoice   string     `json:"invoice" binding:"required"`
	Email     string     `json:"email" binding:"required,email"`
	UserID    *uint      `json:"user_id"` // Nullable for guest users
	Username  string     `json:"username"`
}

type CartItem struct {
	ID       int     `json:"id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,min=0.01"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}


func generateDigest(currency string, merchantEmail string, cartItems []CartItem, totalPrice float64) string {
	// Generate random salt
	salt := fmt.Sprintf("%d", time.Now().UnixNano())

	// Build product info string
	var productParts []string
	for _, item := range cartItems {
		productParts = append(productParts, 
			fmt.Sprintf("%d:%d:%.2f", item.ID, item.Quantity, item.Price))
	}

	// Combine all parts with delimiter
	parts := []string{
		currency,
		merchantEmail,
		salt,
		strings.Join(productParts, "|"),
		fmt.Sprintf("%.2f", totalPrice),
	}
	combined := strings.Join(parts, "||")

	// Generate SHA-256 hash
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func calculateTotal(items []CartItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func buildPayPalItems(items []CartItem) []paypal.Item {
	var paypalItems []paypal.Item
	for _, item := range items {
		paypalItems = append(paypalItems, paypal.Item{
			Name:        item.Name,
			UnitAmount:  &paypal.Money{Currency: "USD", Value: fmt.Sprintf("%.2f", item.Price)},
			Quantity:    fmt.Sprintf("%d", item.Quantity),
			SKU:         fmt.Sprintf("%d", item.ID),
			Description: item.Name,
		})
	}
	return paypalItems
}


type PayPalWebhookEvent struct {
	EventType string `json:"event_type"`
	Resource  struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		CustomID   string `json:"custom_id"`
		CreateTime string `json:"create_time"`
	} `json:"resource"`
}

func PayPalWebhookHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading webhook body: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Parse webhook event
		var event PayPalWebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Log the webhook event
		log.Printf("PayPal Webhook Event: %s, Status: %s", event.EventType, event.Resource.Status)


		// todo
		// Update order status based on event type
		if strings.HasPrefix(event.EventType, "PAYMENT.") {
			// Find order by PayPal transaction ID
			var order models.Order
			if err := db.Where("digest = ?", event.Resource.ID).First(&order).Error; err != nil {
				log.Printf("Order not found for PayPal ID: %s", event.Resource.ID)
			} else {
				// Map PayPal status to our status
				newStatus := map[string]string{
					"COMPLETED": "completed",
					"PENDING":   "pending",
					"FAILED":    "failed",
					"REFUNDED":  "refunded",
				}[event.Resource.Status]

				if newStatus != "" {
					db.Model(&order).Update("status", newStatus)
					log.Printf("Updated order %d status to %s", order.ID, newStatus)
				}
			}
		}
		//todo end

		// Respond with 200 OK
		w.WriteHeader(http.StatusOK)
	}
}

func CheckoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var orderReq PayPalOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Calculate total price
		totalPrice := calculateTotal(orderReq.CartItems)

		// Generate and store order
		order := models.Order{
			Currency:      "USD",
			MerchantEmail: orderReq.Email,
			Salt:         fmt.Sprintf("%d", time.Now().UnixNano()),
			TotalPrice:   totalPrice,
			UserID:       orderReq.UserID,
			Username:     orderReq.Username,
			Digest:       generateDigest("USD", orderReq.Email, orderReq.CartItems, totalPrice),
			CreatedAt:    time.Now(),
		}

		// Add products to order
		for _, item := range orderReq.CartItems {
			order.Products = append(order.Products, models.OrderProduct{
				ProductID: uint(item.ID),
				Quantity:  item.Quantity,
				Price:     item.Price,
			})
		}

		// Save order to database
		if err := db.Create(&order).Error; err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		// Initialize PayPal client
		clientID := os.Getenv("PAYPAL_CLIENT_ID")
		secret := os.Getenv("PAYPAL_SECRET")
		client, err := paypal.NewClient(clientID, secret, paypal.APIBaseSandBox)
		if err != nil {
			http.Error(w, "Failed to initialize PayPal client", http.StatusInternalServerError)
			return
		}

		// Get access token
		_, err = client.GetAccessToken(context.Background())
		if err != nil {
			http.Error(w, "Failed to authenticate with PayPal", http.StatusInternalServerError)
			return
		}

		// Build purchase units
		var purchaseUnits []paypal.PurchaseUnitRequest
		for _, item := range orderReq.CartItems {
			purchaseUnits = append(purchaseUnits, paypal.PurchaseUnitRequest{
				ReferenceID: orderReq.Invoice,
				Description: item.Name,
				Amount: &paypal.PurchaseUnitAmount{
					Currency: "USD",
					Value:    fmt.Sprintf("%.2f", item.Price*float64(item.Quantity)),
					Breakdown: &paypal.PurchaseUnitAmountBreakdown{
						ItemTotal: &paypal.Money{
							Currency: "USD",
							Value:    fmt.Sprintf("%.2f", calculateTotal(orderReq.CartItems)),
						},
					},
				},
				Items: buildPayPalItems(orderReq.CartItems),
			})
		}

		// Create PayPal order
		paypalOrder, err := client.CreateOrder(
			context.Background(),
			paypal.OrderIntentCapture,
			purchaseUnits,
			nil,
			&paypal.ApplicationContext{
				ReturnURL: "http://localhost:3000/",
				CancelURL: "http://localhost:3000/",
			},
		)
		if err != nil {
			http.Error(w, "Failed to create PayPal order", http.StatusInternalServerError)
			return
		}

		// Return approval URL to client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"approvalUrl": paypalOrder.Links[1].Href,
		})
	}
}
