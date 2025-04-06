package models

import (
	"time"
)

type CartItem struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	ListingID     int     `json:"listing_id"`
	BookID        int     `json:"book_id"`
	Price         float32 `json:"price"`
	AllowOffers   bool    `json:"allow_offers"`
	SellerID      int     `json:"seller_id"`
	BookTitle     string  `json:"book_title"`
	BookAuthor    string  `json:"book_author"`
	CoverImageURL string  `json:"cover_image_url"`
	ImageURL      string  `json:"image_url,omitempty"` // Added for first listing image
	Status        string  `json:"status"`

}

// ListingDetails combines Listing + Book Data.
type ListingDetails struct {
	ListingID     int            `json:"listing_id"`
	SellerID      int            `json:"seller_id"`
	BookID        int            `json:"book_id"`
	Price         float32        `json:"price"`
	Status        string         `json:"status"`
	AllowOffers   bool           `json:"allow_offers"`
	SellerNote    string         `json:"seller_note" db:"seller_note"`

	// Book details
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Description   string    `json:"description,omitempty"`
	Language      string    `json:"language,omitempty"`
	ISBN          string    `json:"isbn,omitempty"`
	Publisher     string    `json:"publisher,omitempty"`
	PublishDate   time.Time `json:"publish_date,omitempty"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	AverageRating string    `json:"average_rating,omitempty"`
	NumRatings    string    `json:"num_ratings,omitempty"`
	ImageURLs     []string  `json:"image_urls"`
}

type OfferItem struct {
    ID            int     `json:"id"`
    ListingID     int     `json:"listing_id"`
    BuyerID       int     `json:"buyer_id"`
    OfferedPrice  float64 `json:"offered_price"`
    Status        string  `json:"status"`
    BookID        int     `json:"book_id"`
    BookTitle     string  `json:"book_title"`
    BookAuthor    string  `json:"book_author"`
    CoverImageURL string  `json:"cover_image_url"`
    ImageURL      string  `json:"image_url,omitempty"`
    SellerID      int     `json:"seller_id"`
    BuyerFirstName string  `json:"buyer_first_name"` // New
    BuyerLastName  string  `json:"buyer_last_name"`  // New
    BuyerPicture   string  `json:"buyer_picture_profile"` // New
	InitialPrice   string  `json:"initial_price"`
	Avaibility     string  `json:"avaibility"`
}