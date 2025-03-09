package models

import "time"

type Product struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"catid"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
