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

	"backend/models"

	"github.com/plutov/paypal/v4"
	"gorm.io/gorm"
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

func generateDigest(currency string, merchantEmail string, salt string, cartItems []CartItem, totalPrice float64) string {
	log.Printf("Generating digest with:\nCurrency: %s\nMerchantEmail: %s\nSalt: %s", currency, merchantEmail, salt)

	// Build product info string
	var productParts []string
	for _, item := range cartItems {
		productStr := fmt.Sprintf("%d:%d:%.2f", item.ID, item.Quantity, item.Price)
		productParts = append(productParts, productStr)
		log.Printf("CartItem - ID: %d, Price: %.2f, Quantity: %d", item.ID, item.Price, item.Quantity)
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
	log.Printf("Combined string before hashing: %s", combined)

	// Generate SHA-256 hash
	hash := sha256.Sum256([]byte(combined))
	digest := hex.EncodeToString(hash[:])
	log.Printf("Generated digest: %s", digest)
	return digest
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
			UnitAmount:  &paypal.Money{Currency: "HKD", Value: fmt.Sprintf("%.2f", item.Price)},
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
		ID            string `json:"id"`
		Status        string `json:"status"`
		CustomID      string `json:"custom_id"`
		CreateTime    string `json:"create_time"`
		PurchaseUnits []struct {
			ReferenceID string `json:"reference_id"`
		} `json:"purchase_units"`
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
		log.Println("PayPal Webhook Event: ", event)
		log.Printf("PayPal Webhook Event: %s, Status: %s", event.EventType, event.Resource.Status)

		// todo
		// Update order status when payment is approved
		if event.Resource.Status == "APPROVED" {
			// Get invoice number from purchase_units if available
			invoice := event.Resource.CustomID
			if len(event.Resource.PurchaseUnits) > 0 && event.Resource.PurchaseUnits[0].ReferenceID != "" {
				invoice = event.Resource.PurchaseUnits[0].ReferenceID
			}

			if invoice == "" {
				log.Printf("No invoice number found in webhook event")
				return
			}

			// Find full order details by invoice with manual time parsing
			var order struct {
				models.Order
				CreatedAtStr string `gorm:"column:created_at"`
				UpdatedAtStr string `gorm:"column:updated_at"`
			}
			if err := db.Table("orders").Preload("Products").Where("invoice = ?", invoice).First(&order).Error; err != nil {
				log.Printf("Order not found for invoice: %s", invoice)
				return
			}
			
			// Parse timestamps manually
			createdAt, _ := time.Parse("2006-01-02 15:04:05", order.CreatedAtStr)
			updatedAt, _ := time.Parse("2006-01-02 15:04:05", order.UpdatedAtStr)
			order.Order.CreatedAt = createdAt
			order.Order.UpdatedAt = updatedAt

			// Check if already processed
			if order.Status == "approved" {
				log.Printf("Order %d already approved, skipping", order.ID)
				return
			}

			// Rebuild cart items from order products
			var cartItems []CartItem
			for _, p := range order.Products {
				cartItems = append(cartItems, CartItem{
					ID:       int(p.ProductID),
					Price:    p.Price,
					Quantity: p.Quantity,
				})
			}

			// Regenerate and verify digest
			fmt.Println("order id", order.ID)
			newDigest := generateDigest(order.Currency, order.MerchantEmail, order.Salt, cartItems, order.TotalPrice)
			
			if newDigest != order.Digest {
				log.Printf("Digest mismatch for order %d, possible tampering", order.ID)
				return
			}
			
			// Update status using the original Order model
			if err := db.Model(&order.Order).Where("id = ?", order.ID).Update("status", "approved").Error; err != nil {
				log.Printf("Failed to update order status: %v", err)
				return
			}
			log.Printf("Successfully updated order %d status to approved", order.ID)

			// Save verified order record
			verifiedOrder := models.VerifiedOrder{
				OrderID:    order.ID,
				Invoice:    order.Invoice,
				UserID:     order.UserID,
				Username:   order.Username,
				Email:      order.MerchantEmail,
				TotalPrice: order.TotalPrice,
				Currency:   order.Currency,
				Status:     "approved",
			}

			// Add products to verified order
			for _, p := range order.Products {
				verifiedOrder.Products = append(verifiedOrder.Products, models.VerifiedOrderProduct{
					ProductID: p.ProductID,
					Quantity:  p.Quantity,
					Price:     p.Price,
				})
			}

			// Save to database
			if err := db.Create(&verifiedOrder).Error; err != nil {
				log.Printf("Failed to save verified order: %v", err)
				return
			}
			log.Printf("Successfully saved verified order %d with %d products", verifiedOrder.ID, len(verifiedOrder.Products))
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
		fmt.Println(orderReq.Email)
		fmt.Println(orderReq.Invoice)
		// Generate and store order
		var timesalt = time.Now().UnixNano();
		order := models.Order{
			Currency:      "HKD",
			MerchantEmail: orderReq.Email,
			Salt:          fmt.Sprintf("%d", timesalt),
			TotalPrice:    totalPrice,
			UserID:        orderReq.UserID,
			Username:      orderReq.Username,
			Digest:        generateDigest("HKD", orderReq.Email, fmt.Sprintf("%d", timesalt), orderReq.CartItems, totalPrice),
			Invoice:       orderReq.Invoice,
			CreatedAt:     time.Now(),
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
					Currency: "HKD",
					Value:    fmt.Sprintf("%.2f", item.Price*float64(item.Quantity)),
					Breakdown: &paypal.PurchaseUnitAmountBreakdown{
						ItemTotal: &paypal.Money{
							Currency: "HKD",
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
				ReturnURL: "https://s02.iems5718.ie.cuhk.edu.hk/",
				CancelURL: "https://s02.iems5718.ie.cuhk.edu.hk/",
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
