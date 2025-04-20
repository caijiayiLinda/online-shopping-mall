package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/plutov/paypal/v4"
)

type PayPalOrderRequest struct {
	CartItems []CartItem `json:"cartItems" binding:"required,min=1"`
	Invoice   string     `json:"invoice" binding:"required"`
	Email     string     `json:"email" binding:"required,email"`
}

type CartItem struct {
	ID       int     `json:"id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,min=0.01"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
}


func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var orderReq PayPalOrderRequest
	log.Print(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
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
	order, err := client.CreateOrder(
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
		"approvalUrl": order.Links[1].Href,
	})
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
