package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
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

func (ur *UserRepository) CreateBankAccount(ctx context.Context, bank *models.BankAccount) (int, error) {
	query := `
		INSERT INTO bank_accounts (user_id, bank_name, account_number, account_holder_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`

	result, err := ur.db.ExecContext(ctx, query, bank.UserID, bank.BankName, bank.AccountNumber, bank.AccountHolderName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Println("✅ Bank account created with ID:", id)
	return int(id), nil
}


func (ur *UserRepository) GetGender(ctx context.Context, userID int) (string, error) {
	var gender string
	query := "SELECT gender FROM users WHERE id = ?"

	err := ur.db.QueryRowContext(ctx, query, userID).Scan(
		gender,
	)
	if err == sql.ErrNoRows {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return gender, nil
}

// Updated UserRepository
func (ur *UserRepository) GetAllUserReview(ctx context.Context) ([]models.UserReview, error) {
	query := `
        SELECT id, user_id, book_id, rating
        FROM book_reviews
    `

	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying book_reviews: %w", err)
	}
	defer rows.Close()

	var userReviews []models.UserReview
	for rows.Next() {
		var userReview models.UserReview
		if err := rows.Scan(&userReview.ID, &userReview.UserID, &userReview.BookID, &userReview.Rating); err != nil {
			return nil, fmt.Errorf("error scanning review: %w", err)
		}
		userReviews = append(userReviews, userReview)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}

	return userReviews, nil
}

func (ur *UserRepository) GetAllUserPreferred(ctx context.Context) ([]models.UserPreferred, error) {
	query := `SELECT user_id, genre_id FROM user_preferred_genres`

	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying user_preferred_genres: %w", err)
	}
	defer rows.Close()

	var userPreferreds []models.UserPreferred
	for rows.Next() {
		var userPreferred models.UserPreferred
		if err := rows.Scan(&userPreferred.UserID, &userPreferred.GenreID); err != nil {
			return nil, fmt.Errorf("error scanning review: %w", err)
		}
		userPreferreds = append(userPreferreds, userPreferred)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user_preferred_genres: %w", err)
	}

	return userPreferreds, nil
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

func (ur *UserRepository) AddBookToWishlist(ctx context.Context, userID int, bookID int) (bool, error) {
	var count int
	checkQuery := `SELECT COUNT(*) FROM user_wishlist WHERE user_id = ? AND book_id = ?`
	err := ur.db.QueryRowContext(ctx, checkQuery, userID, bookID).Scan(&count)
	if err != nil {
		return false, err
	}

	// ✅ If the book is already in wishlist, remove it
	if count > 0 {
		deleteQuery := `DELETE FROM user_wishlist WHERE user_id = ? AND book_id = ?`
		_, err = ur.db.ExecContext(ctx, deleteQuery, userID, bookID)
		if err != nil {
			return false, err
		}
		log.Println("✅ Book removed from wishlist for user:", userID, "BookID:", bookID)
		return false, nil // ✅ Book was removed from wishlist
	}

	insertWishlistQuery := `INSERT INTO user_wishlist (user_id, book_id, created_at, updated_at) 
                               VALUES (?, ?, NOW(), NOW())`
	_, err = ur.db.ExecContext(ctx, insertWishlistQuery, userID, bookID)
	if err != nil {
		return false, fmt.Errorf("failed to insert into wishlist: %v", err)
	}
	log.Println("✅ Book added to library for user:", userID, "BookID:", bookID)
	return true, nil
}

func (ur *UserRepository) AddBookToLibrary(ctx context.Context, userID int, bookID int, reading_status int) (bool, error) {

	// Check if the book already exists in user_libraries
	var count int
	checkQuery := `SELECT COUNT(*) FROM user_libraries WHERE user_id = ? AND book_id = ?`
	err := ur.db.QueryRowContext(ctx, checkQuery, userID, bookID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking library: %v", err)
	}

	if count == 0 { // Book not in user_libraries, add it
		insertLibraryQuery := `INSERT INTO user_libraries (user_id, book_id, reading_status, created_at, updated_at) 
                               VALUES (?, ?, ?, NOW(), NOW())`
		_, err = ur.db.ExecContext(ctx, insertLibraryQuery, userID, bookID, reading_status)
		if err != nil {
			return false, fmt.Errorf("failed to insert into library: %v", err)
		}
		log.Println("✅ Book added to library for user:", userID, "BookID:", bookID, "ReadingStatus:", reading_status)
	} else {
		log.Println("ℹ️ Book already exists in library for user:", userID, "BookID:", bookID, "Skipping library insert")
	}
	return false, nil
}

func (ur *UserRepository) AddBookToListing(ctx context.Context, userID int, bookID int, price float32, allowOffer bool, imageURLs []string, sellerNote string) (bool, error) {

	query := `INSERT INTO listings (seller_id, book_id, price, allow_offers, seller_note, created_at, updated_at) 
                  VALUES (?, ?, ?, ?, ?, NOW(), NOW()) 
                  ON DUPLICATE KEY UPDATE price = VALUES(price), allow_offers = VALUES(allow_offers), status = 'for_sale', updated_at = NOW()`

	result, err := ur.db.ExecContext(ctx, query, userID, bookID, price, allowOffer, sellerNote)
	if err != nil {
		return false, fmt.Errorf("failed to insert into listings: %v", err)
	}
	listingID, err := result.LastInsertId()
	if err != nil {
		return false, fmt.Errorf("failed to get listing ID: %v", err)
	}

	// Add images to listing_images
	for _, url := range imageURLs {
		_, err = ur.db.ExecContext(ctx,
			"INSERT INTO listing_images (listing_id, image_url, created_at) VALUES (?, ?, NOW())",
			listingID, url)
		if err != nil {
			return false, fmt.Errorf("failed to add image: %v", err)
		}
	}
	log.Println("✅ Listing added/updated for user:", userID, "BookID:", bookID)

	return false, nil
}

func (ur *UserRepository) FindByID(ctx context.Context, userID int) (*models.GetMe, error) {
	var getMe models.GetMe

	query := `
	SELECT id, email, first_name, last_name, picture_profile, picture_background, 
	phone_number, quote, bio, role, gender, address
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
		&getMe.Gender,
		&getMe.Address,
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

	var bankID int
	err = ur.db.QueryRowContext(ctx, "SELECT id FROM bank_accounts WHERE user_id = ? LIMIT 1", userID).Scan(&bankID)
	if err == sql.ErrNoRows {
		getMe.HasBankAccount = false
	} else if err != nil {
		return nil, err
	} else {
		getMe.HasBankAccount = true
	}
	
	return &getMe, nil
}



func (ur *UserRepository) GetUserPreferredGenres(ctx context.Context, userID int) ([]models.Genre, error) {
	query := `
	SELECT g.id, g.name 
	FROM user_preferred_genres upg
	JOIN genres g ON upg.genre_id = g.id
	WHERE upg.user_id = ?
	`

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferredGenres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		preferredGenres = append(preferredGenres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return preferredGenres, nil
}

func (ur *UserRepository) AddUserPreferredGenres(ctx context.Context, userID int, genreIDs []int) error {
	// Begin transaction
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Delete existing preferences to avoid duplicates
	deleteQuery := `DELETE FROM user_preferred_genres WHERE user_id = ?`
	_, err = tx.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert new preferred genres
	insertQuery := `INSERT INTO user_preferred_genres (user_id, genre_id) VALUES (?, ?)`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, genreID := range genreIDs {
		_, err := stmt.ExecContext(ctx, userID, genreID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateGender(ctx context.Context, userID int, gender string) error {
	query := `UPDATE users SET gender = ? WHERE id = ?`
	result, err := ur.db.ExecContext(ctx, query, gender, userID)
	if err != nil {
		log.Printf("Error updating gender: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}
	log.Printf("Rows affected: %d for userID %d", rowsAffected, userID)
	// Optionally check if user exists, but don't fail if no rows changed
	var exists int
	err = ur.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}
	if exists == 0 {
		return errors.New("no user found with the given ID")
	}
	return nil // Success even if rowsAffected == 0
}

func (ur *UserRepository) SaveProfileImage(userID int, imageURL string) error {
	query := "UPDATE users SET picture_profile = ? WHERE id = ?"
	_, err := ur.db.Exec(query, imageURL, userID)
	if err != nil {
		return err
	}
	return nil
}

// func (ur *UserRepository) SavePostImage(ctx context.Context, postID int, imageURL string) error {
//     query := `
//         INSERT INTO post_images (post_id, image_url)
//         VALUES (?, ?)
//     `
//     _, err := ur.db.ExecContext(ctx, query, postID, imageURL)
//     if err != nil {
//         return fmt.Errorf("failed to save post image: %v", err)
//     }
//     return nil
// }

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

func (ur *UserRepository) EditProfile(ctx context.Context, userID int, first_name string, last_name string, address string, quote string, bio string) error {
	query := `
		UPDATE users
		SET
		first_name = COALESCE(?, first_name),
		last_name = COALESCE(?, last_name),
		address = COALESCE(?, address),
		quote = COALESCE(?, quote), 
		bio = COALESCE(?, bio),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		`
	_, err := ur.db.ExecContext(ctx, query, first_name, last_name, address, quote, bio, userID)
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
	query := `SELECT id, user_id, book_id, reading_status
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
		if err := rows.Scan(&lib.ID, &lib.UserID, &lib.BookID, &lib.ReadingStatus); err != nil {
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

func (ur *UserRepository) GetAllListings(ctx context.Context) ([]models.UserListing, error) {
	query := `SELECT id, seller_id, book_id, price, status, allow_offers, created_at, updated_at 
	          FROM listings 
	          WHERE status = 'for_sale'`

	rows, err := ur.db.QueryContext(ctx, query)
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

	return listings, nil
}

func (ur *UserRepository) GetAllListingsByBookID(ctx context.Context, userID int, bookID int) ([]models.UserListing, error) {
	query := `SELECT id, seller_id, book_id, price, status, allow_offers
	          FROM listings 
	          WHERE book_id = ? AND status = 'for_sale'`

	rows, err := ur.db.QueryContext(ctx, query, bookID)
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

		imageQuery := `
        SELECT image_url 
        FROM listing_images 
        WHERE listing_id = ?
    `
		rows, err := ur.db.QueryContext(ctx, imageQuery, listing.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving listing images: %w", err)
		}
		defer rows.Close()

		var imageURLs []string
		for rows.Next() {
			var imageURL string
			if err := rows.Scan(&imageURL); err != nil {
				return nil, fmt.Errorf("error scanning listing image: %w", err)
			}
			imageURLs = append(imageURLs, imageURL)
		}
		listing.ImageURLs = imageURLs

		listings = append(listings, listing)
	}

	log.Println("Fetched listings for marketplace (excluding user:", userID, ")")
	return listings, nil
}

func (ur *UserRepository) GetMyListings(ctx context.Context, userID int) ([]models.UserListing, error) {
	query := `SELECT id, seller_id, book_id, price, status, allow_offers
	          FROM listings 
	          WHERE seller_id = ? AND status IN ('sold', 'for_sale')`

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
		imageQuery := `
			SELECT image_url 
			FROM listing_images 
			WHERE listing_id = ?
		`
		rows, err := ur.db.QueryContext(ctx, imageQuery, listing.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving listing images: %w", err)
		}
		defer rows.Close()

		var imageURLs []string
		for rows.Next() {
			var imageURL string
			if err := rows.Scan(&imageURL); err != nil {
				return nil, fmt.Errorf("error scanning listing image: %w", err)
			}
			imageURLs = append(imageURLs, imageURL)
		}
		listing.ImageURLs = imageURLs
		listings = append(listings, listing)
	}

	log.Println("Fetched listings for marketplace (excluding user:", userID, ")")
	return listings, nil
}

func (ur *UserRepository) GetPurchasedListingsByUserID(ctx context.Context, userID int) ([]models.MyPurchase, error) {
	query := `
		SELECT 
			l.id AS listing_id,
			l.book_id,
			l.seller_id,
			l.price,
			b.title AS book_title,
			u.first_name,
			u.last_name,
			u.phone_number,
			u.picture_profile,
			t.created_at AS transaction_time
		FROM transactions t
		JOIN listings l ON t.listing_id = l.id
		JOIN users u ON l.seller_id = u.id
		JOIN books b ON l.book_id = b.id
		WHERE t.buyer_id = ? AND t.payment_status = 'completed'
	`

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []models.MyPurchase

	for rows.Next() {
		var p models.MyPurchase
		err := rows.Scan(
			&p.ListingID,
			&p.BookID,
			&p.SellerID,
			&p.Price,
			&p.BookTitle,
			&p.SellerFirstName,
			&p.SellerLastName,
			&p.SellerPhone,
			&p.SellerProfileImg,
			&p.TransactionTime,
		)
		if err != nil {
			return nil, err
		}

		// Get first image from listing_images table
		imgQuery := `SELECT image_url FROM listing_images WHERE listing_id = ? LIMIT 1`
		var imageURL string
		err = ur.db.QueryRowContext(ctx, imgQuery, p.ListingID).Scan(&imageURL)
		if err == nil {
			p.ImageURL = imageURL
		}
		
		purchases = append(purchases, p)
	}

	return purchases, nil
}

func (ur *UserRepository) GetMyOrders(ctx context.Context, sellerID int) ([]models.MyOrder, error) {
	query := `
		SELECT 
			t.listing_id,
			b.title AS book_title,
			t.transaction_amount AS price,
			t.created_at AS transaction_time,
			(
				SELECT image_url 
				FROM listing_images 
				WHERE listing_id = t.listing_id 
				LIMIT 1
			) AS image_url,
			u.id AS buyer_id,
			u.first_name AS buyer_first_name,
			u.last_name AS buyer_last_name,
			u.phone_number AS buyer_phone,
			u.address AS buyer_address,
			u.picture_profile AS buyer_profile_image,
			b.id AS book_id
		FROM transactions t
		JOIN listings l ON t.listing_id = l.id
		JOIN books b ON l.book_id = b.id
		JOIN users u ON t.buyer_id = u.id
		WHERE l.seller_id = ? AND t.payment_status = 'completed'
		ORDER BY t.created_at DESC;
	`

	rows, err := ur.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.MyOrder
	for rows.Next() {
		var o models.MyOrder
		err := rows.Scan(
			&o.ListingID,
			&o.BookTitle,
			&o.Price,
			&o.TransactionTime,
			&o.ImageURL,
			&o.BuyerID,
			&o.BuyerFirstName,
			&o.BuyerLastName,
			&o.BuyerPhone,
			&o.BuyerAddress,
			&o.BuyerProfileImage,
			&o.BookID,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (ur *UserRepository) GetUsersByBookInWishlist(ctx context.Context, bookID int) ([]models.WishlistUser, error) {
	query := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.picture_profile
        FROM user_wishlist uw
        JOIN users u ON uw.user_id = u.id
        WHERE uw.book_id = ?
    `

	rows, err := ur.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.WishlistUser
	for rows.Next() {
		var user models.WishlistUser
		if err := rows.Scan(&user.UserID, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}


func (ur *UserRepository) RemoveListing(ctx context.Context, userID int, listingID int) error {
	// Verify the user owns the listing and it's not already sold or removed
	query := `
        UPDATE listings 
        SET status = 'removed', updated_at = NOW()
        WHERE id = ? AND seller_id = ? AND status NOT IN ('sold', 'removed')
    `
	result, err := ur.db.ExecContext(ctx, query, listingID, userID)
	if err != nil {
		return fmt.Errorf("error removing listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found, already sold/removed, or not owned by user")
	}

	log.Printf("Listing %d marked as removed for user %d", listingID, userID)
	return nil
}

func (ur *UserRepository) IsBookInWishlist(ctx context.Context, userID int, bookID int) (bool, error) {
	// Query to check if the book exists in the wishlist
	query := `SELECT COUNT(*) 
	          FROM user_wishlist 
	          WHERE user_id = ? AND book_id = ?`

	var count int
	// Execute the query and scan the result
	err := ur.db.QueryRowContext(ctx, query, userID, bookID).Scan(&count)
	if err != nil {
		return false, err
	}

	// Return true if the book is in the wishlist, false otherwise
	return count > 0, nil
}

func (ur *UserRepository) GetWishlistByUserID(ctx context.Context, userID int) ([]models.Book, error) {
	query := `SELECT b.id, b.title, b.author, b.cover_image_url, b.description, COALESCE(br.average_rating, 0) AS average_rating
	          FROM user_wishlist ul
	          JOIN books b ON ul.book_id = b.id
              LEFT JOIN book_ratings br ON b.id = br.book_id
	          WHERE ul.user_id = ?`
              

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlist []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.CoverImageURL, &book.Description, &book.AverageRating)
		if err != nil {
			return nil, err
		}
		wishlist = append(wishlist, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlist, nil
}

// GetListingByID retrieves a listing along with seller's Omise account ID
func (ur *UserRepository) GetListingByID(ctx context.Context, listingID int) (*models.ListingDetails, error) {
	query := `
        SELECT 
            l.id AS listing_id, l.seller_id, l.book_id, l.price, l.status, l.allow_offers, l.seller_note,
            b.title, b.author, b.description, b.language, b.isbn, b.publisher, 
            b.publish_date, b.cover_image_url, 
            COALESCE(br.average_rating, 0) AS average_rating, 
            COALESCE(br.num_ratings, 0) AS num_ratings
        FROM listings l
        JOIN books b ON l.book_id = b.id
        LEFT JOIN book_ratings br ON b.id = br.book_id
        JOIN users u ON l.seller_id = u.id
        WHERE l.id = ?
        LIMIT 1
    `

	var listing models.ListingDetails
	err := ur.db.QueryRowContext(ctx, query, listingID).Scan(
		&listing.ListingID, &listing.SellerID, &listing.BookID,
		&listing.Price, &listing.Status, &listing.AllowOffers, &listing.SellerNote,
		&listing.Title, &listing.Author, &listing.Description,
		&listing.Language, &listing.ISBN, &listing.Publisher,
		&listing.PublishDate, &listing.CoverImageURL,
		&listing.AverageRating, &listing.NumRatings, // Fetch seller's Omise account ID
	)

	if err != nil {
		return nil, err
	}

	imageQuery := `
        SELECT image_url 
        FROM listing_images 
        WHERE listing_id = ?
    `
	rows, err := ur.db.QueryContext(ctx, imageQuery, listingID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving listing images: %w", err)
	}
	defer rows.Close()

	var imageURLs []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, fmt.Errorf("error scanning listing image: %w", err)
		}
		imageURLs = append(imageURLs, imageURL)
	}
	listing.ImageURLs = imageURLs

	return &listing, nil
}

// AddToCart adds a listing to a user's cart
func (ur *UserRepository) AddToCart(ctx context.Context, userID int, listingID int) (int, error) {
	// First, verify that the listing exists and is available
	var status string
	checkQuery := `
        SELECT status 
        FROM listings 
        WHERE id = ? AND status = 'for_sale'
    `
	err := ur.db.QueryRowContext(ctx, checkQuery, listingID).Scan(&status)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("listing %d is not available or does not exist", listingID)
	}
	if err != nil {
		return 0, fmt.Errorf("error checking listing availability: %w", err)
	}

	// Check if the listing is already in the user's cart
	var existingCount int
	countQuery := `
        SELECT COUNT(*) 
        FROM cart 
        WHERE user_id = ? AND listing_id = ?
    `
	err = ur.db.QueryRowContext(ctx, countQuery, userID, listingID).Scan(&existingCount)
	if err != nil {
		return 0, fmt.Errorf("error checking existing cart entry: %w", err)
	}

	if existingCount > 0 {
		return 0, fmt.Errorf("listing %d is already in user %d's cart", listingID, userID)
	}

	// Insert the listing into the cart
	insertQuery := `
        INSERT INTO cart (user_id, listing_id, created_at)
        VALUES (?, ?, NOW())
    `
	result, err := ur.db.ExecContext(ctx, insertQuery, userID, listingID)
	if err != nil {
		return 0, fmt.Errorf("failed to add listing to cart: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get cart entry ID: %w", err)
	}

	log.Printf("Successfully added listing %d to user %d's cart with cart ID %d", listingID, userID, id)
	return int(id), nil
}

// GetCart retrieves all listings in a user's cart
func (ur *UserRepository) GetCart(ctx context.Context, userID int) ([]models.CartItem, error) {
	query := `
        SELECT 
            c.id,
            c.user_id,
            c.listing_id,
            l.book_id,
            l.price,
            l.allow_offers,
            l.seller_id,
            b.title,
            b.author,
            b.cover_image_url,
            (SELECT image_url 
             FROM listing_images li 
             WHERE li.listing_id = l.id 
             ORDER BY li.id ASC 
             LIMIT 1) AS image_url,
            l.status  -- Added to retrieve listing status
        FROM cart c
        JOIN listings l ON c.listing_id = l.id
        JOIN books b ON l.book_id = b.id
        WHERE c.user_id = ? AND l.status IN ('for_sale', 'sold', 'reserved')  -- Updated to include 'for_sale' and 'sold'
    `

	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying cart: %w", err)
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ListingID,
			&item.BookID,
			&item.Price,
			&item.AllowOffers,
			&item.SellerID,
			&item.BookTitle,
			&item.BookAuthor,
			&item.CoverImageURL,
			&item.ImageURL, // First image
			&item.Status,   // Added for status
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart item: %w", err)
		}
		cartItems = append(cartItems, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	return cartItems, nil
}

// RemoveFromCart removes a listing from a user's cart
func (ur *UserRepository) RemoveFromCart(ctx context.Context, userID int, listingID int) error {
	query := `
        DELETE FROM cart 
        WHERE user_id = ? AND listing_id = ?
    `
	result, err := ur.db.ExecContext(ctx, query, userID, listingID)
	if err != nil {
		return fmt.Errorf("failed to remove from cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no cart entry found for user %d and listing %d", userID, listingID)
	}

	log.Printf("Removed listing %d from user %d's cart", listingID, userID)
	return nil
}

// AddToOffers adds an offer for a listing by a buyer
func (ur *UserRepository) AddToOffers(ctx context.Context, buyerID int, listingID int, offeredPrice float64) (int, error) {
	// Verify that the listing exists and allows offers
	var allowOffers bool
	var status string
	checkQuery := `
        SELECT allow_offers, status
        FROM listings
        WHERE id = ? AND status = 'for_sale'
    `
	err := ur.db.QueryRowContext(ctx, checkQuery, listingID).Scan(&allowOffers, &status)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("listing %d is not available or does not exist", listingID)
	}
	if err != nil {
		return 0, fmt.Errorf("error checking listing: %w", err)
	}
	if !allowOffers {
		return 0, fmt.Errorf("listing %d does not allow offers", listingID)
	}

	// Check if the buyer already has a pending/accepted offer for this listing
	var existingCount int
	countQuery := `
        SELECT COUNT(*)
        FROM offers
        WHERE buyer_id = ? AND listing_id = ? AND status IN ('pending', 'accepted')
    `
	err = ur.db.QueryRowContext(ctx, countQuery, buyerID, listingID).Scan(&existingCount)
	if err != nil {
		return 0, fmt.Errorf("error checking existing offer: %w", err)
	}
	if existingCount > 0 {
		return 0, fmt.Errorf("buyer %d already has an active offer for listing %d", buyerID, listingID)
	}

	// Insert the offer
	insertQuery := `
        INSERT INTO offers (listing_id, buyer_id, offered_price, status, created_at)
        VALUES (?, ?, ?, 'pending', NOW())
    `
	result, err := ur.db.ExecContext(ctx, insertQuery, listingID, buyerID, offeredPrice)
	if err != nil {
		return 0, fmt.Errorf("failed to add offer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get offer ID: %w", err)
	}

	log.Printf("Successfully added offer for listing %d by buyer %d with offer ID %d", listingID, buyerID, id)
	return int(id), nil
}

// GetOffers retrieves all offers made by a buyer
func (ur *UserRepository) GetBuyerOffers(ctx context.Context, buyerID int) ([]models.OfferItem, error) {
	query := `
        SELECT 
            o.id,
            o.listing_id,
            o.buyer_id,
            o.offered_price,
            o.status,
            l.book_id,
            b.title,
            b.author,
            b.cover_image_url,
            (SELECT image_url 
             FROM listing_images li 
             WHERE li.listing_id = l.id 
             ORDER BY li.id ASC 
             LIMIT 1) AS image_url,
            l.seller_id,
			l.price AS initial_price,
			l.status AS avaibility
        FROM offers o
        JOIN listings l ON o.listing_id = l.id
        JOIN books b ON l.book_id = b.id
        WHERE o.buyer_id = ?
    `
	rows, err := ur.db.QueryContext(ctx, query, buyerID)
	if err != nil {
		return nil, fmt.Errorf("error querying offers: %w", err)
	}
	defer rows.Close()

	var offers []models.OfferItem
	for rows.Next() {
		var item models.OfferItem
		err := rows.Scan(
			&item.ID,
			&item.ListingID,
			&item.BuyerID,
			&item.OfferedPrice,
			&item.Status,
			&item.BookID,
			&item.BookTitle,
			&item.BookAuthor,
			&item.CoverImageURL,
			&item.ImageURL,
			&item.SellerID,
			&item.InitialPrice,
			&item.Avaibility,

		)
		if err != nil {
			return nil, fmt.Errorf("error scanning offer item: %w", err)
		}
		offers = append(offers, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating offers: %w", err)
	}

	return offers, nil
}

// repository/user_repository.go
func (ur *UserRepository) GetSellerOffers(ctx context.Context, sellerID int) ([]models.OfferItem, error) {
	query := `
        SELECT 
            o.id,
            o.listing_id,
            o.buyer_id,
            o.offered_price,
            o.status,
            l.book_id,
            b.title,
            b.author,
            b.cover_image_url,
            (SELECT image_url 
             FROM listing_images li 
             WHERE li.listing_id = l.id 
             ORDER BY li.id ASC 
             LIMIT 1) AS image_url,
            l.seller_id,
            u.first_name AS buyer_first_name,
            u.last_name AS buyer_last_name,
            u.picture_profile AS buyer_picture_profile,
			l.price AS initial_price
        FROM offers o
        JOIN listings l ON o.listing_id = l.id
        JOIN books b ON l.book_id = b.id
        JOIN users u ON o.buyer_id = u.id
        WHERE l.seller_id = ?
    `
	rows, err := ur.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("error querying seller offers: %w", err)
	}
	defer rows.Close()

	var offers []models.OfferItem
	for rows.Next() {
		var item models.OfferItem
		err := rows.Scan(
			&item.ID,
			&item.ListingID,
			&item.BuyerID,
			&item.OfferedPrice,
			&item.Status,
			&item.BookID,
			&item.BookTitle,
			&item.BookAuthor,
			&item.CoverImageURL,
			&item.ImageURL,
			&item.SellerID,
			&item.BuyerFirstName,
			&item.BuyerLastName,
			&item.BuyerPicture,
			&item.InitialPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning seller offer item: %w", err)
		}
		offers = append(offers, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating seller offers: %w", err)
	}

	return offers, nil
}

// RemoveFromOffers removes an offer (e.g., buyer retracts it)
func (ur *UserRepository) RemoveFromOffers(ctx context.Context, buyerID int, listingID int) error {
	query := `
        DELETE FROM offers
        WHERE buyer_id = ? AND listing_id = ? AND status = 'pending'
    `
	result, err := ur.db.ExecContext(ctx, query, buyerID, listingID)
	if err != nil {
		return fmt.Errorf("failed to remove offer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no pending offer found for buyer %d and listing %d", buyerID, listingID)
	}

	log.Printf("Removed offer for listing %d by buyer %d", listingID, buyerID)
	return nil
}

// AcceptOffer: Seller accepts an offer

// repository/user_repository.go
func (ur *UserRepository) AcceptOffer(ctx context.Context, sellerID int, offerID int) error {
	query := `
        UPDATE offers o
        JOIN listings l ON o.listing_id = l.id
        SET o.status = 'accepted',
            o.updated_at = NOW()
        WHERE o.id = ? 
        AND l.seller_id = ? 
        AND o.status = 'pending'
        AND l.status = 'for_sale'
    `
	result, err := ur.db.ExecContext(ctx, query, offerID, sellerID)
	if err != nil {
		return fmt.Errorf("failed to accept offer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no pending offer found with ID %d for seller %d or listing not for sale", offerID, sellerID)
	}

	log.Printf("Seller %d accepted offer %d", sellerID, offerID)
	return nil
}

// RejectOffer: Seller rejects an offer
func (ur *UserRepository) RejectOffer(ctx context.Context, sellerID int, offerID int) error {
	query := `
        UPDATE offers o
        JOIN listings l ON o.listing_id = l.id
        SET o.status = 'rejected',
            o.updated_at = NOW()
        WHERE o.id = ? 
        AND l.seller_id = ? 
        AND o.status = 'pending'
    `
	result, err := ur.db.ExecContext(ctx, query, offerID, sellerID)
	if err != nil {
		return fmt.Errorf("failed to reject offer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no pending offer found with ID %d for seller %d", offerID, sellerID)
	}

	log.Printf("Seller %d rejected offer %d", sellerID, offerID)
	return nil
}

// repository/user_repository.go
func (ur *UserRepository) GetAcceptedOffer(ctx context.Context, offerID int) (*models.OfferItem, error) {
	query := `
        SELECT 
            o.id,
            o.listing_id,
            o.buyer_id,
            o.offered_price,
            o.status,
            l.book_id,
            b.title,
            b.author,
            b.cover_image_url,
            (SELECT image_url 
             FROM listing_images li 
             WHERE li.listing_id = l.id 
             ORDER BY li.id ASC 
             LIMIT 1) AS image_url,
            l.seller_id,
            u.first_name AS buyer_first_name,
            u.last_name AS buyer_last_name,
            u.picture_profile AS buyer_picture_profile
        FROM offers o
        JOIN listings l ON o.listing_id = l.id
        JOIN books b ON l.book_id = b.id
        JOIN users u ON o.buyer_id = u.id
        WHERE o.id = ? AND o.status = 'accepted'
    `
	var item models.OfferItem
	err := ur.db.QueryRowContext(ctx, query, offerID).Scan(
		&item.ID,
		&item.ListingID,
		&item.BuyerID,
		&item.OfferedPrice,
		&item.Status,
		&item.BookID,
		&item.BookTitle,
		&item.BookAuthor,
		&item.CoverImageURL,
		&item.ImageURL,
		&item.SellerID,
		&item.BuyerFirstName,
		&item.BuyerLastName,
		&item.BuyerPicture,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no accepted offer found with ID %d", offerID)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching accepted offer: %w", err)
	}
	return &item, nil
}

func (ur *UserRepository) GetOfferByID(ctx context.Context, offerID int) (*models.OfferItem, error) {
	query := `
        SELECT 
            o.id,
            o.listing_id,
            o.buyer_id,
            o.offered_price,
            o.status,
            l.book_id,
            b.title,
            b.author,
            b.cover_image_url,
            (SELECT image_url 
             FROM listing_images li 
             WHERE li.listing_id = l.id 
             ORDER BY li.id ASC 
             LIMIT 1) AS image_url,
            l.seller_id,
            u.first_name AS buyer_first_name,
            u.last_name AS buyer_last_name,
            u.picture_profile AS buyer_picture_profile
        FROM offers o
        JOIN listings l ON o.listing_id = l.id
        JOIN books b ON l.book_id = b.id
        JOIN users u ON o.buyer_id = u.id
        WHERE o.id = ?
    `
	var item models.OfferItem
	err := ur.db.QueryRowContext(ctx, query, offerID).Scan(
		&item.ID,
		&item.ListingID,
		&item.BuyerID,
		&item.OfferedPrice,
		&item.Status,
		&item.BookID,
		&item.BookTitle,
		&item.BookAuthor,
		&item.CoverImageURL,
		&item.ImageURL,
		&item.SellerID,
		&item.BuyerFirstName,
		&item.BuyerLastName,
		&item.BuyerPicture,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no offer found with ID %d", offerID)
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching offer: %w", err)
	}
	return &item, nil
}

// repository/user_repository.go
// repository/user_repository.go
func (ur *UserRepository) ReserveListingForOffer(ctx context.Context, listingID int, buyerID int) (bool, error) {
	query := `
        UPDATE listings l
        JOIN offers o ON o.listing_id = l.id
        SET l.status = 'reserved',
            l.reserved_expires_at = NOW() + INTERVAL 2 MINUTE,
            l.updated_at = NOW()
        WHERE l.id = ?
        AND o.buyer_id = ?
        AND o.status = 'accepted'
        AND l.status = 'for_sale'
    `
	result, err := ur.db.ExecContext(ctx, query, listingID, buyerID)
	if err != nil {
		return false, fmt.Errorf("failed to reserve listing for offer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return false, nil // Listing wasn’t reserved (e.g., not for_sale or offer not accepted)
	}

	log.Printf("Listing %d reserved for buyer %d via offer, expires at %s", listingID, buyerID, time.Now().Add(2*time.Minute))
	return true, nil
}

// repository/user_repository.go
// repository/user_repository.go
func (ur *UserRepository) RevertOfferReservation(ctx context.Context, listingID int, offerID int) error {
	query := `
        UPDATE listings 
        SET status = 'for_sale', 
            reserved_expires_at = NULL, 
            updated_at = NOW()
        WHERE id = ? 
        AND status = 'reserved'
    `
	result, err := ur.db.ExecContext(ctx, query, listingID)
	if err != nil {
		return fmt.Errorf("failed to revert listing: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no reserved listing found with ID %d", listingID)
	}

	// Offer stays 'accepted'—buyer can try again
	return nil
}

// ReserveListing reserves a listing atomically with a timeout
// repository/user_repository.go
func (ur *UserRepository) ReserveListing(ctx context.Context, listingID int, buyerID int) (bool, error) {
	query := `
        UPDATE listings 
        SET status = 'reserved',
            reserved_expires_at = NOW() + INTERVAL 2 MINUTE,
            updated_at = NOW()
        WHERE id = ?
        AND status = 'for_sale'
    `
	result, err := ur.db.ExecContext(ctx, query, listingID)
	if err != nil {
		return false, fmt.Errorf("failed to reserve listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return false, nil // Listing wasn’t reserved (e.g., not for_sale)
	}

	log.Printf("Listing %d reserved for buyer %d, expires at %s", listingID, buyerID, time.Now().Add(2*time.Minute))
	return true, nil
}

// MarkListingAsSold updates the listing status to sold
// repository/user_repository.go
// repository/user_repository.go
func (ur *UserRepository) MarkListingAsSold(ctx context.Context, listingID, buyerID int) error {
	tx, err := ur.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Step 1: Mark the listing as sold
	queryListing := `
        UPDATE listings 
        SET status = 'sold',
            reserved_expires_at = NULL,
            updated_at = NOW()
        WHERE id = ?
        AND status = 'reserved'
    `
	result, err := tx.ExecContext(ctx, queryListing, listingID)
	if err != nil {
		return fmt.Errorf("failed to mark listing as sold: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected for listing: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("listing %d is not in reserved state", listingID)
	}

	// Step 2: Handle offers
	// - If this is an offer payment, keep the winning offer 'accepted' (or update to 'completed')
	// - Reject all other offers
	queryOffers := `
        UPDATE offers 
        SET status = 'rejected',
            updated_at = NOW()
        WHERE listing_id = ?
        AND buyer_id != ?
        AND status IN ('pending', 'accepted')
    `
	_, err = tx.ExecContext(ctx, queryOffers, listingID, buyerID)
	if err != nil {
		return fmt.Errorf("failed to reject other offers: %w", err)
	}

	// // Optional: Update the winning offer to 'completed' if it’s an offer payment
	// queryWinningOffer := `
    //     UPDATE offers 
    //     SET status = 'completed',
    //         updated_at = NOW()
    //     WHERE listing_id = ?
    //     AND buyer_id = ?
    //     AND status = 'accepted'
    // `
	// _, err = tx.ExecContext(ctx, queryWinningOffer, listingID, buyerID)
	// if err != nil {
	// 	return fmt.Errorf("failed to update winning offer: %w", err)
	// }

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Listing %d marked as sold for buyer %d, offers updated", listingID, buyerID)
	return nil
}

// ExpireReservedListing reverts a listing to for_sale if payment isn’t completed
// repository/user_repository.go
func (ur *UserRepository) ExpireReservedListing(ctx context.Context, listingID int) error {
	query := `
        UPDATE listings 
        SET status = 'for_sale',
            reserved_expires_at = NULL,
            updated_at = NOW()
        WHERE id = ?
        AND status = 'reserved'
        AND reserved_expires_at <= NOW()
    `
	result, err := ur.db.ExecContext(ctx, query, listingID)
	if err != nil {
		return fmt.Errorf("failed to expire listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no expired reserved listing found with ID %d", listingID)
	}

	log.Printf("Listing %d reservation expired", listingID)
	return nil
}

// CreateTransaction records a new transaction
func (ur *UserRepository) CreateTransaction(ctx context.Context, stripe_session_id string, buyerID int, listingID int, offer_id *int, amount float64, status string) error {

	query := `INSERT INTO transactions (stripe_session_id, buyer_id, listing_id, offer_id, transaction_amount, payment_status, created_at, updated_at)
             VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	
	var offerValue interface{}
	if offer_id != nil {
		offerValue = *offer_id
	} else {
		offerValue = nil
	}

	_, err := ur.db.ExecContext(ctx, query,
		stripe_session_id,
		buyerID,
		listingID,
		offerValue,
		amount,
		status,
	)


	return err
}

// UpdateTransactionStatus updates the transaction status
func (ur *UserRepository) UpdateTransactionStatus(ctx context.Context, listingID int, status string) error {
	query := `UPDATE transactions SET payment_status = ?, updated_at = NOW() 
	          WHERE listing_id = ? AND payment_status = 'reserved'`
	_, err := ur.db.ExecContext(ctx, query, status, listingID)
	return err
}



// CleanupExpiredListings runs as a background job
// repository/user_repository.go
func (ur *UserRepository) CleanupExpiredListings(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Less frequent since webhook handles most cases
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			query := `
                SELECT id 
                FROM listings 
                WHERE status = 'reserved' 
                AND reserved_expires_at IS NOT NULL 
                AND reserved_expires_at <= NOW()
            `
			rows, err := ur.db.QueryContext(ctx, query)
			if err != nil {
				log.Println("❌ Cleanup Error:", err)
				continue
			}
			var listingIDs []int
			for rows.Next() {
				var id int
				if err := rows.Scan(&id); err != nil {
					log.Println("❌ Scan Error:", err)
					continue
				}
				listingIDs = append(listingIDs, id)
			}
			if err := rows.Close(); err != nil {
				log.Println("❌ Rows Close Error:", err)
			}
			for _, id := range listingIDs {
				if err := ur.ExpireReservedListing(ctx, id); err != nil {
					log.Println("❌ Expire Error for listing", id, ":", err)
				} else {
					log.Printf("Cleanup: Expired listing %d reverted to for_sale", id)
				}
			}
		case <-ctx.Done():
			log.Println("CleanupExpiredListings stopped")
			return
		}
	}
}

// IsListingReserved checks if a listing is reserved and still active
// repository/user_repository.go
func (ur *UserRepository) IsListingReserved(ctx context.Context, listingID int) (bool, bool, error) {
	var status string
	var reservedExpiresAt *time.Time
	query := `
        SELECT status, reserved_expires_at 
        FROM listings 
        WHERE id = ?
    `
	err := ur.db.QueryRowContext(ctx, query, listingID).Scan(&status, &reservedExpiresAt)
	if err != nil {
		return false, false, fmt.Errorf("failed to check listing status: %w", err)
	}
	if status != "reserved" {
		return false, false, nil
	}
	if reservedExpiresAt == nil || time.Now().Before(*reservedExpiresAt) {
		return true, false, nil // Reserved and not expired
	}
	return true, true, nil // Reserved but expired
}

// repository/user_repository.go
func (ur *UserRepository) CreatePost(ctx context.Context, userID int, content string, imageURLs []string, genreID *int, bookID *int) (models.Post, error) {
	// Insert post
	query := `
        INSERT INTO posts (user_id, content, genre_id, book_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, NOW(), NOW())
    `
	result, err := ur.db.ExecContext(ctx, query, userID, content, genreID, bookID)
	if err != nil {
		return models.Post{}, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return models.Post{}, err
	}

	// Insert image URLs only if provided
	if len(imageURLs) > 0 {
		for _, url := range imageURLs {
			_, err := ur.db.ExecContext(ctx, "INSERT INTO post_images (post_id, image_url) VALUES (?, ?)", postID, url)
			if err != nil {
				return models.Post{}, err
			}
		}
	}

	post := models.Post{
		ID:        int(postID),
		UserID:    userID,
		Content:   content,
		GenreID:   genreID,
		BookID:    bookID,
		ImageURLs: imageURLs,
		CreatedAt: time.Now(),
	}
	return post, nil
}

func (ur *UserRepository) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	// Fetch all posts
	query := `
        SELECT id, user_id, content, genre_id, book_id, created_at, updated_at
        FROM posts ORDER BY created_at DESC
    `
	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var genreID sql.NullInt64
		var bookID sql.NullInt64
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &genreID, &bookID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		if genreID.Valid {
			id := int(genreID.Int64)
			post.GenreID = &id
		}
		if bookID.Valid {
			id := int(bookID.Int64)
			post.BookID = &id
		}

		// Fetch image URLs for this post
		imageRows, err := ur.db.QueryContext(ctx, "SELECT image_url FROM post_images WHERE post_id = ?", post.ID)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()

		var imageURLs []string
		for imageRows.Next() {
			var url string
			if err := imageRows.Scan(&url); err != nil {
				return nil, err
			}
			imageURLs = append(imageURLs, url)
		}
		post.ImageURLs = imageURLs

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// Optional: Fetch post with images
func (ur *UserRepository) GetPost(ctx context.Context, id int) (models.Post, error) {
	var post models.Post

	// Fetch post details
	row := ur.db.QueryRowContext(ctx, "SELECT id, user_id, content, created_at FROM posts WHERE id = ?", id)
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, sql.ErrNoRows // Let caller handle "not found"
		}
		return models.Post{}, err
	}

	// Fetch image URLs
	rows, err := ur.db.QueryContext(ctx, "SELECT image_url FROM post_images WHERE post_id = ?", id)
	if err != nil {
		return models.Post{}, err
	}
	defer rows.Close()

	var imageURLs []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return models.Post{}, err
		}
		imageURLs = append(imageURLs, url)
	}
	if err = rows.Err(); err != nil {
		return models.Post{}, err
	}

	post.ImageURLs = imageURLs
	return post, nil
}

// // GetAllPosts fetches all posts with their image URLs
// func (ur *UserRepository) GetAllPosts(ctx context.Context) ([]models.Post, error) {
// 	// Fetch all posts
// 	rows, err := ur.db.QueryContext(ctx, "SELECT id, user_id, content, created_at FROM posts ORDER BY created_at DESC")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch posts: %v", err)
// 	}
// 	defer rows.Close()

// 	var posts []models.Post
// 	for rows.Next() {
// 		var post models.Post
// 		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt); err != nil {
// 			return nil, fmt.Errorf("failed to scan post: %v", err)
// 		}

// 		// Fetch image URLs for this post
// 		imageRows, err := ur.db.QueryContext(ctx, "SELECT image_url FROM post_images WHERE post_id = ?", post.ID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to fetch images for post %d: %v", post.ID, err)
// 		}
// 		defer imageRows.Close()

// 		var imageURLs []string
// 		for imageRows.Next() {
// 			var url string
// 			if err := imageRows.Scan(&url); err != nil {
// 				return nil, fmt.Errorf("failed to scan image URL: %v", err)
// 			}
// 			imageURLs = append(imageURLs, url)
// 		}
// 		if err = imageRows.Err(); err != nil {
// 			return nil, err
// 		}
// 		post.ImageURLs = imageURLs
// 		posts = append(posts, post)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return posts, nil
// }

// GetPostsByUserID fetches all posts by a specific user with their image URLs
// repository/user_repository.go
func (ur *UserRepository) GetPostsByUserID(ctx context.Context, userID int) ([]models.Post, error) {
	query := `
        SELECT id, user_id, content, genre_id, book_id, created_at, updated_at
        FROM posts
        WHERE user_id = ?
        ORDER BY created_at DESC
    `
	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts for user %d: %v", userID, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var genreID sql.NullInt64
		var bookID sql.NullInt64
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &genreID, &bookID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %v", err)
		}
		if genreID.Valid {
			id := int(genreID.Int64)
			post.GenreID = &id
		}
		if bookID.Valid {
			id := int(bookID.Int64)
			post.BookID = &id
		}

		// Fetch image URLs
		imageRows, err := ur.db.QueryContext(ctx, "SELECT image_url FROM post_images WHERE post_id = ?", post.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch images for post %d: %v", post.ID, err)
		}
		defer imageRows.Close()

		var imageURLs []string
		for imageRows.Next() {
			var url string
			if err := imageRows.Scan(&url); err != nil {
				return nil, fmt.Errorf("failed to scan image URL: %v", err)
			}
			imageURLs = append(imageURLs, url)
		}
		if err = imageRows.Err(); err != nil {
			return nil, err
		}
		post.ImageURLs = imageURLs

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil // Return empty slice if no posts, not nil
}

// GetPostByPostID fetches a single post by its ID with image URLs
func (ur *UserRepository) GetPostByPostID(ctx context.Context, postID int) (models.Post, error) {
	var post models.Post
	row := ur.db.QueryRowContext(ctx, "SELECT id, user_id, content, created_at FROM posts WHERE id = ?", postID)
	err := row.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, fmt.Errorf("post %d not found", postID)
		}
		return models.Post{}, fmt.Errorf("failed to fetch post %d: %v", postID, err)
	}

	// Fetch image URLs
	imageRows, err := ur.db.QueryContext(ctx, "SELECT image_url FROM post_images WHERE post_id = ?", postID)
	if err != nil {
		return models.Post{}, fmt.Errorf("failed to fetch images for post %d: %v", postID, err)
	}
	defer imageRows.Close()

	var imageURLs []string
	for imageRows.Next() {
		var url string
		if err := imageRows.Scan(&url); err != nil {
			return models.Post{}, fmt.Errorf("failed to scan image URL: %v", err)
		}
		imageURLs = append(imageURLs, url)
	}
	if err = imageRows.Err(); err != nil {
		return models.Post{}, err
	}
	post.ImageURLs = imageURLs

	return post, nil
}

// CreateComment adds a new comment to a post
func (ur *UserRepository) CreateComment(ctx context.Context, postID, userID int, content string) (models.Comment, error) {
	query := "INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, NOW())"
	result, err := ur.db.ExecContext(ctx, query, postID, userID, content)
	if err != nil {
		return models.Comment{}, fmt.Errorf("failed to create comment: %v", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, err
	}

	var createdAt time.Time
	err = ur.db.QueryRowContext(ctx, "SELECT created_at FROM comments WHERE id = ?", commentID).Scan(&createdAt)
	if err != nil {
		return models.Comment{}, err
	}

	comment := models.Comment{
		ID:        int(commentID),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: createdAt,
	}
	return comment, nil
}

// GetCommentsByPostID fetches all comments for a post
func (ur *UserRepository) GetCommentsByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	rows, err := ur.db.QueryContext(ctx, "SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at ASC", postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %v", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %v", err)
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateLike adds a like to a post
func (ur *UserRepository) CreateLike(ctx context.Context, postID, userID int) (models.Like, error) {
	query := "INSERT INTO post_likes (post_id, user_id, created_at) VALUES (?, ?, NOW())"
	result, err := ur.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return models.Like{}, fmt.Errorf("user %d already liked post %d", userID, postID)
		}
		return models.Like{}, fmt.Errorf("failed to create like: %v", err)
	}

	likeID, err := result.LastInsertId()
	if err != nil {
		return models.Like{}, err
	}

	var createdAt time.Time
	err = ur.db.QueryRowContext(ctx, "SELECT created_at FROM post_likes WHERE id = ?", likeID).Scan(&createdAt)
	if err != nil {
		return models.Like{}, err
	}

	like := models.Like{
		ID:        int(likeID),
		PostID:    postID,
		UserID:    userID,
		CreatedAt: createdAt,
	}
	return like, nil
}

// RemoveLike removes a like from a post
func (ur *UserRepository) RemoveLike(ctx context.Context, postID, userID int) error {
	query := "DELETE FROM post_likes WHERE post_id = ? AND user_id = ?"
	result, err := ur.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove like: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no like found for post %d by user %d", postID, userID)
	}
	return nil
}

// GetLikeCountByPostID returns the number of likes for a post
func (ur *UserRepository) GetLikeCountByPostID(ctx context.Context, postID int) (int, error) {
	var count int
	err := ur.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM post_likes WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get like count: %v", err)
	}
	return count, nil
}

// IsPostLikedByUser checks if a user has liked a post
func (ur *UserRepository) IsPostLikedByUser(ctx context.Context, postID, userID int) (bool, error) {
	var exists int
	err := ur.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check like status: %v", err)
	}
	return exists > 0, nil
}

