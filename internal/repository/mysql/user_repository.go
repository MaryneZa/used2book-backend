package mysql

import (
	"context"
	"database/sql"
	"log"
	"used2book-backend/internal/models"
	"fmt"
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


// GetAllUsers retrieves all users from the database
func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]models.GetAllUsers, error) {
    query := `
    SELECT 
        id, email, first_name, last_name, picture_profile, picture_background,
        phone_number, gender, quote, bio, role
    FROM users
    `

    // Execute the query
    rows, err := ur.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("error querying users: %w", err)
    }
    defer rows.Close()

    // Slice to hold the results
    var users []models.GetAllUsers

    // Iterate through the result set
    for rows.Next() {
        var user models.GetAllUsers
        err := rows.Scan(
            &user.ID, &user.Email, &user.FirstName, &user.LastName,
            &user.ProfilePicture, &user.BackgroundPicture, &user.PhoneNumber,
            &user.Gender, &user.Quote, &user.Bio, &user.Role,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning user: %w", err)
        }
        users = append(users, user)
    }

    // Check for errors during iteration
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating users: %w", err)
    }

    return users, nil
}


func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, hashed_password, COALESCE(first_name, '') AS first_name, COALESCE(last_name, '') AS last_name , COALESCE(picture_profile, '') AS picture_profile, COALESCE(picture_background, '') AS picture_background FROM users WHERE email = ?"

	err := ur.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.ProfilePicture, &user.BackgroundPicture,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	query := `INSERT INTO users (first_name, last_name, email, hashed_password, provider, role, created_at, updated_at) VALUES (?,?, ?, ?, ?, ?, ?, ?)`

	result, err := ur.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.HashedPassword, user.Provider, user.Role, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}
	// user.ID = int(id)
	log.Println("id: ", id)
	return int(id), nil
}

func (ur *UserRepository) AddBookToLibrary(ctx context.Context, userID int, bookID int, ownStatus string, price float32, allowOffer bool) error {
	query := `INSERT INTO user_libraries (user_id, book_id, status, created_at, updated_at) 
	          VALUES (?, ?, ?, NOW(), NOW()) 
	          ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = NOW()`

	_, err := ur.db.ExecContext(ctx, query, userID, bookID, ownStatus)
	if err != nil {
		return err
	}

	// If the book is "owned", add/update the listing
	if ownStatus == "owned" {
		query = `INSERT INTO listings (seller_id, book_id, price, allow_offers, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, NOW(), NOW()) 
		         ON DUPLICATE KEY UPDATE price = VALUES(price), allow_offers = VALUES(allow_offers), status = 'for_sale', updated_at = NOW()`

		_, err = ur.db.ExecContext(ctx, query, userID, bookID, price, allowOffer)
		if err != nil {
			return err
		}
	}

	log.Println("Book added to library for user:", userID, "BookID:", bookID, "Status:", ownStatus)
	return nil
}


func (ur *UserRepository) FindByID(ctx context.Context, userID int) (*models.GetMe, error) {
	var getMe models.GetMe

	query := `
	SELECT id, email, first_name, last_name, picture_profile, picture_background, 
	phone_number, quote, bio, role 
	FROM users 
	WHERE id = ?
	`

	err := ur.db.QueryRowContext(ctx, query, userID).Scan(
		&getMe.ID,
		&getMe.Email,
		&getMe.FirstName,
		&getMe.LastName,
		&getMe.ProfilePicture,
		&getMe.BackgroundPicture,
		&getMe.PhoneNumber,
		&getMe.Quote,
		&getMe.Bio,
		&getMe.Role,
	)
	if err != nil {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &getMe, nil
}

func (ur *UserRepository) SaveProfileImage(userID int, imageURL string) error {
	query := "UPDATE users SET picture_profile = ? WHERE id = ?"
	_, err := ur.db.Exec(query, imageURL, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) SaveBackgroundImage(userID int, imageURL string) error {
	query := "UPDATE users SET picture_background = ? WHERE id = ?"
	_, err := ur.db.Exec(query, imageURL, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) EditAccountInfo(ctx context.Context, userID int, firstName string, lastName string, phoneNumber sql.NullString) error {
	query := `
		UPDATE users
		SET 
		first_name = COALESCE(?, first_name), 
		last_name = COALESCE(?, last_name),
		phone_number = COALESCE(?, phone_number),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
	_, err := ur.db.ExecContext(ctx, query, firstName, lastName, phoneNumber, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) EditQuote(ctx context.Context, userID int, quote string) error {
	query := `
		UPDATE users
		SET 
		quote = COALESCE(?, quote), 
		WHERE id = ?
		`
	_, err := ur.db.ExecContext(ctx, query, quote, userID)
	if err != nil {
		return err
	}
	return nil
}
func (ur *UserRepository) EditBio(ctx context.Context, userID int, bio string) error {
	query := `
	UPDATE users
	SET 
	bio = COALESCE(?, bio), 
	WHERE id = ?
	`
	_, err := ur.db.ExecContext(ctx, query, bio, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) EditPhoneNumber(ctx context.Context, userID int, phoneNumber string) error {
	query := `
	UPDATE users
	SET 
	phone_number = COALESCE(?, phone_number),
	updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`
	_, err := ur.db.ExecContext(ctx, query, phoneNumber, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) IsPhoneNumberTaken(ctx context.Context, phoneNumber string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE phone_number = ?"

	err := ur.db.QueryRowContext(ctx, query, phoneNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ur *UserRepository) EditName(ctx context.Context, userID int, firstName string, lastName string) error {
	query := `
		UPDATE users
		SET 
		first_name = COALESCE(?, first_name), 
		last_name = COALESCE(?, last_name),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
	_, err := ur.db.ExecContext(ctx, query, firstName, lastName, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) EditPreferrence(ctx context.Context, userID int, quote string, bio string) error {
	query := `
		UPDATE users
		SET 
		quote = COALESCE(?, quote), 
		bio = COALESCE(?, bio),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
	_, err := ur.db.ExecContext(ctx, query, quote, bio, userID)
	if err != nil {
		return err
	}
	return nil
}

// CountBooks checks how many books exist in the database
func (ur *UserRepository) CountUsers() (int, error) {
	var count int
	err := ur.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ur *UserRepository) GetUserLibrary(ctx context.Context, userID int) ([]models.UserLibrary, error) {
	query := `SELECT id, user_id, book_id, status
	          FROM user_libraries 
	          WHERE user_id = ?`

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libraries []models.UserLibrary

	for rows.Next() {
		var lib models.UserLibrary
		if err := rows.Scan(&lib.ID, &lib.UserID, &lib.BookID, &lib.Status); err != nil {
			return nil, err
		}
		libraries = append(libraries, lib)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	log.Println("User Library fetched for user:", userID)
	return libraries, nil
}


func (ur *UserRepository) GetAllListings(ctx context.Context, userID int) ([]models.UserListing, error) {
	query := `SELECT id, seller_id, book_id, price, status, allow_offers, created_at, updated_at 
	          FROM listings 
	          WHERE seller_id != ? AND status = 'for_sale'`

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []models.UserListing
	for rows.Next() {
		var listing models.UserListing
		if err := rows.Scan(&listing.ID, &listing.SellerID, &listing.BookID, &listing.Price, &listing.Status, &listing.AllowOffer); err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	log.Println("Fetched listings for marketplace (excluding user:", userID, ")")
	return listings, nil
}


// func (ur *UserRepository) AddUserLibraryEntry(ctx context.Context, entry models.UserLibrary) (int, error) {
// 	query := `
// 	INSERT INTO user_libraries (user_id, book_id, status, personal_notes, favorite_quotes)
// 	VALUES (?, ?, ?, ?, ?)
// 	`

// 	result, err := ur.db.ExecContext(ctx, query, entry.UserID, entry.BookID, entry.Status, entry.PersonalNotes, entry.FavoriteQuotes)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to insert user library entry: %w", err)
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
// 	}

// 	return int(id), nil
// }

// // GetUserLibrary retrieves a user's library based on their user ID
// func (ur *UserRepository) GetUserLibrary(ctx context.Context, userID int) ([]models.UserLibrary, error) {
// 	query := `
// 	SELECT id, user_id, book_id, status, personal_notes, favorite_quotes, created_at, updated_at
// 	FROM user_libraries WHERE user_id = ?
// 	`

// 	rows, err := ur.db.QueryContext(ctx, query, userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query user library: %w", err)
// 	}
// 	defer rows.Close()

// 	var library []models.UserLibrary
// 	for rows.Next() {
// 		var entry models.UserLibrary
// 		if err := rows.Scan(&entry.ID, &entry.UserID, &entry.BookID, &entry.Status, &entry.PersonalNotes, &entry.FavoriteQuotes, &entry.CreatedAt, &entry.UpdatedAt); err != nil {
// 			return nil, fmt.Errorf("failed to scan user library row: %w", err)
// 		}
// 		library = append(library, entry)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating through user library rows: %w", err)
// 	}

// 	return library, nil
// }

// // UpdateUserLibrary updates an existing user library entry
// func (ur *UserRepository) UpdateUserLibrary(ctx context.Context, entry models.UserLibrary) error {
// 	query := `
// 	UPDATE user_libraries
// 	SET status = ?, personal_notes = ?, favorite_quotes = ?, updated_at = CURRENT_TIMESTAMP
// 	WHERE id = ? AND user_id = ?
// 	`

// 	_, err := ur.db.ExecContext(ctx, query, entry.Status, entry.PersonalNotes, entry.FavoriteQuotes, entry.ID, entry.UserID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update user library entry: %w", err)
// 	}

// 	return nil
// }

// // DeleteUserLibraryEntry removes an entry from the user_libraries table
// func (ur *UserRepository) DeleteUserLibraryEntry(ctx context.Context, id int, userID int) error {
// 	query := `DELETE FROM user_libraries WHERE id = ? AND user_id = ?`
// 	_, err := ur.db.ExecContext(ctx, query, id, userID)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user library entry: %w", err)
// 	}

// 	return nil
// }



