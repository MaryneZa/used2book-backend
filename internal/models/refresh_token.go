package models

import "time"

type RefreshToken struct {
    ID        int       `db:"id" json:"id"`
    UserID    int       `db:"user_id" json:"user_id"`
    Token     string    `db:"token" json:"token"`
    DeviceInfo     string    `db:"device_info" json:"device_info"`
    ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
