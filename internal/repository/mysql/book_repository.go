package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"used2book-backend/internal/models"
)

// BookRepository struct
type BookRepository struct {
	db *sql.DB
}

// NewBookRepository initializes a new BookRepository
func NewBookRepository(db *sql.DB) *BookRepository {
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	return &BookRepository{db}
}

// InsertBook inserts a book into MySQL
func (br *BookRepository) InsertBook(ctx context.Context, book models.Book) (int, error) {
	query := `
	INSERT INTO books (title, author, description, language, isbn, 
	publisher, publish_date, cover_image_url) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the insert query
	result, err := br.db.ExecContext(ctx, query,
		book.Title, book.Author, book.Description, book.Language, book.ISBN,
		book.Publisher, book.PublishDate, book.CoverImageURL,
	)
	if err != nil {
		return 0, err
	}

	// Retrieve the last inserted ID
	bookID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert int64 to int
	finalBookID := int(bookID)

	// Insert an initial rating entry for the book
	_, err = br.db.ExecContext(ctx, "INSERT INTO book_ratings (book_id) VALUES (?)", finalBookID)
	if err != nil {
		return 0, err
	}

	return finalBookID, nil
}

// GetAllBooks retrieves all books from the database
func (br *BookRepository) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	query := `
    SELECT id, title, author, description, language, isbn, 
           publisher, publish_date, cover_image_url 
    FROM books
    `

	// Execute the query
	rows, err := br.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying books: %w", err)
	}
	defer rows.Close()

	// Slice to hold the results
	var books []models.Book

	// Iterate through the result set
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.Description,
			&book.Language, &book.ISBN, &book.Publisher,
			&book.PublishDate, &book.CoverImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning book: %w", err)
		}
		books = append(books, book)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating books: %w", err)
	}

	return books, nil
}


// GetBookByID retrieves a book by its ID, including its rating information
func (br *BookRepository) GetBookByID(ctx context.Context, bookID int) (*models.Book, error) {
	query := `
    SELECT b.id, b.title, b.author, b.description, b.language, b.isbn, 
           b.publisher, b.publish_date, b.cover_image_url,
           COALESCE(br.average_rating, 0) AS average_rating, 
           COALESCE(br.num_ratings, 0) AS num_ratings
    FROM books b
    LEFT JOIN book_ratings br ON b.id = br.book_id
    WHERE b.id = ?
    `

	// Execute the query
	row := br.db.QueryRowContext(ctx, query, bookID)

	// Variable to hold the result
	var book models.Book

	// Scan the result into the book struct
	err := row.Scan(
		&book.ID, &book.Title, &book.Author, &book.Description,
		&book.Language, &book.ISBN, &book.Publisher,
		&book.PublishDate, &book.CoverImageURL,
		&book.AverageRating, &book.NumRatings, // New fields for rating
	)

	if err == sql.ErrNoRows {
		return nil, nil // No book found with the given ID
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving book by ID: %w", err)
	}

	return &book, nil
}

// GetBookWithRatings fetches a book along with its rating
// func (br *BookRepository) GetBookWithRatings(ctx context.Context, bookID int) (*models.BookWithRatings, error) {
// 	query := `
// 	SELECT b.id, b.title, b.author, b.description, b.language, b.isbn, 
// 		    b.publisher, b.publish_date, b.cover_image_url, 
// 		       br.average_rating, br.num_ratings
// 		FROM books b
// 		LEFT JOIN book_ratings br ON b.id = br.book_id
// 		WHERE b.id = ?
// 	`

// 	var bookWithRatings models.BookWithRatings
// 	err := br.db.QueryRowContext(ctx, query, bookID).Scan(
// 		&bookWithRatings.Book.ID, &bookWithRatings.Book.Title,
// 		&bookWithRatings.Book.Author, &bookWithRatings.Book.Description, &bookWithRatings.Book.Language,
// 		&bookWithRatings.Book.ISBN,
// 		&bookWithRatings.Book.Publisher, &bookWithRatings.Book.PublishDate,
// 		&bookWithRatings.Book.CoverImageURL, &bookWithRatings.Ratings.AverageRating,
// 		&bookWithRatings.Ratings.NumRatings,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &bookWithRatings, nil
// }

// // UpdateUserBookRating allows a user to rate or edit their book rating
// func (br *BookRepository) UpdateUserBookRating(ctx context.Context, userID int, bookID int, newRating float64) error {
// 	var existingRating float64
// 	err := br.db.QueryRowContext(ctx, "SELECT rating FROM user_book_ratings WHERE user_id = ? AND book_id = ?", userID, bookID).Scan(&existingRating)

// 	if err == sql.ErrNoRows {
// 		// Insert new rating
// 		_, err = br.db.ExecContext(ctx, "INSERT INTO user_book_ratings (user_id, book_id, rating) VALUES (?, ?, ?)", userID, bookID, newRating)
// 	} else if err == nil {
// 		// Update existing rating
// 		_, err = br.db.ExecContext(ctx, "UPDATE user_book_ratings SET rating = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND book_id = ?", newRating, userID, bookID)
// 	} else {
// 		return err
// 	}

// 	// Recalculate the average rating
// 	_, err = br.db.ExecContext(ctx, `
// 		UPDATE book_ratings 
// 		SET average_rating = (SELECT COALESCE(AVG(rating), 0) FROM user_book_ratings WHERE book_id = ?),
// 		    num_ratings = (SELECT COUNT(*) FROM user_book_ratings WHERE book_id = ?)
// 		WHERE book_id = ?
// 	`, bookID, bookID, bookID)

// 	return err
// }

// // DeleteUserRating removes a user's rating for a book
// func (br *BookRepository) DeleteUserRating(ctx context.Context, userID int, bookID int) error {
// 	// Delete user's rating from `user_book_ratings`
// 	_, err := br.db.ExecContext(ctx, "DELETE FROM user_book_ratings WHERE user_id = ? AND book_id = ?", userID, bookID)
// 	if err != nil {
// 		return err
// 	}

// 	// Recalculate average rating and number of ratings
// 	_, err = br.db.ExecContext(ctx, `
// 		UPDATE book_ratings 
// 		SET average_rating = (SELECT COALESCE(AVG(rating), 0) FROM user_book_ratings WHERE book_id = ?),
// 		    num_ratings = (SELECT COUNT(*) FROM user_book_ratings WHERE book_id = ?)
// 		WHERE book_id = ?
// 	`, bookID, bookID, bookID)

// 	return err
// }

// CountBooks checks how many books exist in the database
func (br *BookRepository) CountBooks() (int, error) {
	var count int
	err := br.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (br *BookRepository) NumberOfBooks(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM books"

	// Execute the query and scan the result into the count variable
	err := br.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err // Return 0 and the error if the query fails
	}

	return count, nil // Return the count and no error
}

func (br *BookRepository) GetReviewsByBookID(ctx context.Context, bookID int) ([]models.BookReview, error) {
	query := `SELECT br.id, br.user_id, u.first_name, u.last_name, u.picture_profile, br.rating, br.comment, br.created_at, br.updated_at
			  FROM book_reviews br
			  JOIN users u ON br.user_id = u.id
			  WHERE br.book_id = ? 
			  ORDER BY br.created_at DESC`

	rows, err := br.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.BookReview
	for rows.Next() {
		var review models.BookReview
		err := rows.Scan(&review.ID, &review.UserID, &review.FirstName, &review.LastName, &review.UserProfile,&review.Rating, &review.Comment, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}


// Helper functions for parsing values
func parseInt(value string) int {
	var i int
	fmt.Sscanf(value, "%d", &i)
	return i
}

func parseFloat(value string) float64 {
	var f float64
	fmt.Sscanf(value, "%f", &f)
	return f
}

func parseDate(dateStr string) (time.Time, error) {
	cleaned := strings.TrimSpace(dateStr)
	return time.Parse("01/02/06", cleaned)
}

// FetchBooksFromGoogleSheets retrieves book data from Google Sheets
func (br *BookRepository) FetchBooksFromGoogleSheets(sheetID, apiKey string) ([][]string, error) {
	sheetName := "books_1.Best_Books_Ever" // Change this if your sheet name is different
	url := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s", sheetID, sheetName, apiKey)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from Google Sheets: %v", err)
	}
	defer resp.Body.Close()

	// Google Sheets API returns JSON, so we decode it properly
	var result struct {
		Values [][]string `json:"values"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON from Google Sheets: %v", err)
	}

	if len(result.Values) < 2 {
		return nil, fmt.Errorf("Google Sheet is empty or only contains headers")
	}

	return result.Values, nil
}

func (br *BookRepository) SyncBooksFromGoogleSheets(sheetID, apiKey string) error {
	records, err := br.FetchBooksFromGoogleSheets(sheetID, apiKey)
	if err != nil {
		return err
	}

	if len(records) < 2 {
		log.Println("⚠️ No data found in Google Sheets, skipping sync")
		return nil
	}

	ctx := context.Background()
	for i, row := range records[1:30] { // Skip header row
		if len(row) < 18 {
			log.Printf("⚠️ Skipping row %d: insufficient data\n", i+1)
			continue
		}

		// Parse publish date
		publishDate, err := parseDate(strings.TrimSpace(row[14]))
		if err != nil {
			log.Printf("⚠️ Skipping row %d: invalid publish date for book '%s'\n", i+1, row[1])
			log.Println(err)
			continue
		}

		genreList := strings.Split(row[8], ",")
		for i := range genreList {
			// Trim whitespace and any extraneous characters like brackets or quotes
			genreList[i] = strings.TrimSpace(genreList[i])     // Remove surrounding spaces
			genreList[i] = strings.Trim(genreList[i], "[]'\"") // Remove brackets and quotes
		}

		// Insert book details
		book := models.Book{
			Title:       row[1], // title
			Author:      row[3], // author
			Description: row[5], // description
			Language:    row[6], // language
			ISBN:        row[7], // isbn
			// Genres:        row[8],   // genres
			Publisher:     row[13], // publisher
			PublishDate:   publishDate,
			CoverImageURL: row[21], // coverImg
		}

		bookID, err := br.InsertBook(ctx, book)
		if err != nil {
			log.Printf("⚠️ Skipping row %d: Error inserting book '%s' → %v\n", i+1, row[1], err)
			continue
		}

		// Insert or fetch genres and associate with book
		for _, genreName := range genreList {
			genreID, err := br.GetOrInsertGenre(ctx, genreName)
			if err != nil {
				log.Printf("⚠️ Error associating genre '%s' with book '%s': %v\n", genreName, book.Title, err)
				continue
			}

			err = br.AssociateBookWithGenre(ctx, bookID, genreID)
			if err != nil {
				log.Printf("⚠️ Error creating book-genre association for '%s': %v\n", genreName, err)
			}
		}
	}

	log.Println("✅ Books successfully synced from Google Sheets!")
	return nil
}

// GetOrInsertGenre retrieves the genre ID if it exists, or inserts it if it doesn't.
func (br *BookRepository) GetOrInsertGenre(ctx context.Context, genreName string) (int, error) {
	var genreID int
	query := "SELECT id FROM genres WHERE name = ?"
	err := br.db.QueryRowContext(ctx, query, genreName).Scan(&genreID)

	if err == sql.ErrNoRows {
		// Insert new genre
		insertQuery := "INSERT INTO genres (name) VALUES (?)"
		result, err := br.db.ExecContext(ctx, insertQuery, genreName)
		if err != nil {
			return 0, err
		}
		genreID64, _ := result.LastInsertId()
		return int(genreID64), nil
	}

	return genreID, err
}

// AssociateBookWithGenre creates an association between a book and a genre.
func (br *BookRepository) AssociateBookWithGenre(ctx context.Context, bookID, genreID int) error {
	query := "INSERT INTO book_genres (book_id, genre_id) VALUES (?, ?)"
	_, err := br.db.ExecContext(ctx, query, bookID, genreID)
	return err
}

// GetGenresByBookID fetches the genres associated with a given book ID.
func (br *BookRepository) GetGenresByBookID(ctx context.Context, bookID int) ([]string, error) {
	query := `
		SELECT g.name 
		FROM genres g
		JOIN book_genres bg ON g.id = bg.genre_id
		WHERE bg.book_id = ?
	`
	rows, err := br.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []string
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}


func (br *BookRepository) AddBookReview(ctx context.Context, userID int, bookID int, rating float32, comment string) error {
	tx, err := br.db.BeginTx(ctx, nil) // ✅ Use a transaction to ensure atomicity
	if err != nil {
		return err
	}

	// ✅ Step 1: Insert or update the review
	reviewQuery := `INSERT INTO book_reviews (user_id, book_id, rating, comment, created_at, updated_at) 
	                VALUES (?, ?, ?, ?, NOW(), NOW()) 
	                ON DUPLICATE KEY UPDATE rating = VALUES(rating), comment = VALUES(comment), updated_at = NOW()`
	_, err = tx.ExecContext(ctx, reviewQuery, userID, bookID, rating, comment)
	if err != nil {
		tx.Rollback()
		return err
	}

	// ✅ Step 2: Recalculate average rating and number of ratings
	ratingQuery := `
		UPDATE book_ratings 
		SET average_rating = (SELECT COALESCE(AVG(rating), 0) FROM book_reviews WHERE book_id = ?),
		    num_ratings = (SELECT COUNT(*) FROM book_reviews WHERE book_id = ?)
		WHERE book_id = ?`
	_, err = tx.ExecContext(ctx, ratingQuery, bookID, bookID, bookID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// ✅ Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (br *BookRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	query := `SELECT id, name FROM genres ORDER BY name ASC`

	rows, err := br.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}
