package models

import "time"

// Book represents static book metadata.
type Book struct {
    ID            int       `json:"id" db:"id"`
    Title         string    `json:"title" db:"title"`
    Author        string    `json:"author" db:"author"`
    Description   string    `json:"description,omitempty" db:"description"`
    Language      string    `json:"language,omitempty" db:"language"`
    ISBN          string    `json:"isbn,omitempty" db:"isbn"`
    Publisher     string    `json:"publisher,omitempty" db:"publisher"`
    PublishDate   time.Time `json:"publish_date,omitempty" db:"publish_date"`
    CoverImageURL string    `json:"cover_image_url,omitempty" db:"cover_image_url"`
    NumRatings    string    `json:"num_ratings,omitempty" db:"num_ratings"`
    AverageRating string    `json:"average_rating,omitempty" db:"average_rating"`
    CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// BookRatings represents ratings and popularity metrics.
type BookRatings struct {
    ID            int       `json:"id" db:"id"`
    BookID        int       `json:"book_id" db:"book_id"`
    AverageRating float64   `json:"average_rating,omitempty" db:"average_rating"`
    NumRatings    int       `json:"num_ratings,omitempty" db:"num_ratings"`
    CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// BookWithRatings is a combination of Book and its Ratings
type BookWithRatings struct {
    Book    Book        `json:"book"`
    Ratings BookRatings `json:"ratings"`
}

type BookReview struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
    BookID    int       `json:"book_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
    UserProfile string  `json:"picture_profile"`
	Rating    float32   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type AddBookReview struct {
	ID      int     `json:"id,omitempty"` // Auto-generated, no need in request
	UserID  int     `json:"user_id,omitempty"` // Not needed in request (comes from context)
	BookID  int     `json:"book_id"`
	Rating  float32 `json:"rating"`
	Comment string  `json:"comment"`
}

type Genre struct {
	ID      int     `json:"id,omitempty" db:"id"` // Auto-generated, no need in request
	Name      string    `json:"name,omitempty" db:"name"` // Auto-generated, no need in request
}

type BookGenre struct {
	BookID  int     `json:"book_id"`
    GenreID int     `json:"genre_id"`
}



