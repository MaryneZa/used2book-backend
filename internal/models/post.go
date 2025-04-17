package models

import (
	"time"
)



type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    GenreID   *int      `json:"genre_id,omitempty"`
    BookID    *int      `json:"book_id,omitempty"`
    ImageURLs []string  `json:"image_urls,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}



type Comment struct {
    ID             int       `json:"id"`
    PostID         int       `json:"post_id"`
    UserID         int       `json:"user_id"`
    Content        string    `json:"content"`
    CreatedAt      time.Time `json:"created_at"`
    FirstName      string    `json:"first_name"`
    LastName       string    `json:"last_name"`
    PictureProfile string    `json:"picture_profile"`
}

type Like struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}


