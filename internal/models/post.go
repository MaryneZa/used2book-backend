package models

import (
	"time"
)


type Post struct {
    ID        int      `json:"id"`
    UserID    int      `json:"user_id"`
    Content   string   `json:"content"`
    ImageURLs []string `json:"image_urls,omitempty"` // Populated from post_images
    CreatedAt time.Time   `json:"created_at"`
}

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type Like struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}


