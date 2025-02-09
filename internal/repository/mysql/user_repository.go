package mysql

import (
	"context"
	"database/sql"
	"log"
	"used2book-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	if db == nil {
		log.Fatal("database connection is nil")
	}
	return &UserRepository{db}
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, hashed_password, COALESCE(username, '') AS username, COALESCE(picture, '') AS picture FROM users WHERE email = ?"

	err := ur.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.HashedPassword, &user.Name, &user.ProfilePicture,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, hashed_password, provider, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := ur.db.ExecContext(ctx, query, user.Email, user.HashedPassword, user.Provider, user.Role, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}
	// user.ID = int(id)
	log.Println("id: ", id)
	return nil
}

func (ur *UserRepository) FindByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, username, picture FROM users WHERE id = ?"


	err := ur.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.Name, &user.ProfilePicture,
	)
	log.Println("user hey: ", user, " error: ", err)
	

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
