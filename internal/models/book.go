package models

import "time"

// Book represents a book entity.
type Book struct {
    ID            int       `json:"id" db:"id"`
    Title         string    `json:"title" db:"title"`
    Author        string    `json:"author" db:"author"`
    Description   string    `json:"description,omitempty" db:"description"`
    Category      string    `json:"category,omitempty" db:"category"`
    CoverImageURL string    `json:"cover_image_url,omitempty" db:"cover_image_url"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
