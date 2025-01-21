package mysql

import (
	"context"
	"database/sql"
	"used2book-backend/internal/models"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error){
	var user models.User
	query := "SELECT email FROM users WHERE email = ?"
	err := ur.db.QueryRowContext(ctx, query, email).Scan(
		&user.Email,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *models.User) error{
	query := `INSERT INTO users (email, verified_email, name, picture, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := ur.db.ExecContext(ctx, query, user.Email, user.VerifiedEmail, user.Name, user.ProfilePictureURL, user.Role, user.CreatedAt, user.UpdatedAt)
	
	if err != nil {
		return err
	}
	id, err:= result.LastInsertId()

	if err != nil{
		return err
	}
	// user.ID = int(id)
	log.Println("id: ", id)
    return nil
}