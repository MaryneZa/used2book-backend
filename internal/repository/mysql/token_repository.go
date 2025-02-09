package mysql

import (
	"context"
	"database/sql"
	"log"
	"time"
	"used2book-backend/internal/models"
	"used2book-backend/internal/utils"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	if db == nil {
		log.Fatal("database connection is nil")
	}
	return &TokenRepository{db}
}

func (tr *TokenRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, username, picture FROM users WHERE email = ?"

	err := tr.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.ProfilePicture,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (tr *TokenRepository) StoreRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := tr.db.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func (tr *TokenRepository) UpdateRefreshToken(ctx context.Context, userID int, newRefreshToken string) error {
	updateQuery := `UPDATE refresh_tokens SET token = ?, expires_at = ? WHERE user_id = ?`
	_, err := tr.db.ExecContext(ctx, updateQuery, newRefreshToken, utils.RefreshTokenExpiration(), userID)
	return err
}

func (tr *TokenRepository) ValidateRefreshToken(ctx context.Context, token string) (int, error) {
	var userID int
	query := `SELECT user_id FROM refresh_tokens WHERE token = ? AND expires_at > NOW()`
	err := tr.db.QueryRowContext(ctx, query, token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

func (tr *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := tr.db.ExecContext(ctx, query, token)
	return err
}

func (tr *TokenRepository) DeleteRefreshTokenById(ctx context.Context, user_id int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := tr.db.ExecContext(ctx, query, user_id)
	return err
}

func (tr *TokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at <= NOW()`
	_, err := tr.db.ExecContext(ctx, query)
	return err
}

func (tr *TokenRepository) DeleteTokensById(ctx context.Context, userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := tr.db.ExecContext(ctx, query, userID)
	log.Printf("delete token success : %s", err)
	return err
}
