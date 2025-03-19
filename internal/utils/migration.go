package utils

import (
	"log"
)

// RunMigrations executes a series of SQL statements to create tables and constraints.
func RunMigrations() {
	db := GetDB() // Obtain a connected *sql.DB instance.
	if db == nil {
		log.Fatalf("Database connection failed")
	}

	queries := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,

            first_name VARCHAR(255) DEFAULT '',
            last_name VARCHAR(255) DEFAULT '',

            provider ENUM('google','local') NOT NULL,
            
            hashed_password VARCHAR(255),
            phone_number VARCHAR(20) DEFAULT NULL UNIQUE,

            picture_profile VARCHAR(255) DEFAULT '',
            picture_background VARCHAR(255) DEFAULT '',

            gender ENUM('male', 'female', 'other') NOT NULL DEFAULT 'other',

            quote VARCHAR(100) DEFAULT '',
            bio VARCHAR(500) DEFAULT '',

            omise_account_id VARCHAR(255) DEFAULT NULL, 

            role ENUM('user','admin') DEFAULT 'user',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );`,

		// provider_id VARCHAR(255),

		// Books table
		`CREATE TABLE IF NOT EXISTS books (
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255) NOT NULL,
            description TEXT,
            language VARCHAR(50),
            isbn VARCHAR(20) UNIQUE,
            publisher VARCHAR(255),
            publish_date DATE,
            cover_image_url VARCHAR(500),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );`,

		// Create book_ratings table (stores calculated ratings)
		`CREATE TABLE IF NOT EXISTS book_ratings (
            id INT AUTO_INCREMENT PRIMARY KEY,
            book_id INT NOT NULL UNIQUE,
            average_rating DECIMAL(3,2) DEFAULT 0.0,
            num_ratings INT DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
        );`,

		// User Libraries table
		`CREATE TABLE IF NOT EXISTS user_libraries (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            book_id INT NOT NULL,
            status ENUM('owned','not_own','wishlist') NOT NULL,
            personal_notes VARCHAR(255) DEFAULT '' ,
            favorite_quotes VARCHAR(255) DEFAULT '',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

		`CREATE TABLE IF NOT EXISTS listings (
            id INT AUTO_INCREMENT PRIMARY KEY,
            seller_id INT NOT NULL,
            book_id INT NOT NULL,
            price FLOAT NOT NULL,
            status ENUM('for_sale', 'reserved_sale', 'sold', 'removed') DEFAULT 'for_sale',
            reserved_expires_at TIMESTAMP NULL DEFAULT NULL,
            allow_offers BOOLEAN DEFAULT FALSE,
            seller_omise_id VARCHAR(255) DEFAULT NULL,
            seller_note TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

		`CREATE TABLE IF NOT EXISTS listing_images (
            id INT AUTO_INCREMENT PRIMARY KEY,
            listing_id INT NOT NULL,
            image_url VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE
        );`,

		`CREATE TABLE IF NOT EXISTS cart (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            listing_id INT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (listing_id) REFERENCES listings(id)
        );`,

		// Offers table
		`CREATE TABLE IF NOT EXISTS offers (
            id INT AUTO_INCREMENT PRIMARY KEY,
            listing_id INT NOT NULL,
            buyer_id INT NOT NULL,
            offered_price DECIMAL(10,2) NOT NULL,
            status ENUM('pending','accepted','rejected','countered') DEFAULT 'pending',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (listing_id) REFERENCES listings(id),
            FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE
        );`,

		// Transactions table
		`CREATE TABLE IF NOT EXISTS transactions (
            id INT AUTO_INCREMENT PRIMARY KEY,
            buyer_id INT NOT NULL,
            seller_id INT NOT NULL,
            listing_id INT NOT NULL,
            transaction_amount DECIMAL(10,2) NOT NULL,
            payment_status ENUM('reserved','completed','failed') DEFAULT 'reserved',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (listing_id) REFERENCES listings(id)
        );`,

		// Book Reviews table
		`CREATE TABLE IF NOT EXISTS book_reviews (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            book_id INT NOT NULL,
            rating DECIMAL(10,2) NOT NULL,
            comment TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
            );`,

		`CREATE TABLE IF NOT EXISTS post_images (
            id INT AUTO_INCREMENT PRIMARY KEY,
            post_id INT NOT NULL,
            image_url VARCHAR(255) NOT NULL,
            FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
        );`,

		// Comments table
		`CREATE TABLE IF NOT EXISTS comments (
            id INT AUTO_INCREMENT PRIMARY KEY,
            post_id INT NOT NULL,
            user_id INT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
            );`,

		`CREATE TABLE IF NOT EXISTS post_likes (
                id INT AUTO_INCREMENT PRIMARY KEY,
                post_id INT NOT NULL,
                user_id INT NOT NULL,
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                UNIQUE KEY unique_like (post_id, user_id) -- Ensures one like per user per post
            );`,

		// Notifications table
		`CREATE TABLE IF NOT EXISTS notifications (
                id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            message VARCHAR(255),
            is_read BOOLEAN DEFAULT false,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
            );`,

		// genres table
		`CREATE TABLE IF NOT EXISTS genres (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
            );`,

		// Book genres table (Pivot)
		`CREATE TABLE IF NOT EXISTS book_genres (
                id INT AUTO_INCREMENT PRIMARY KEY,
                book_id INT NOT NULL,
            genre_id INT NOT NULL,
            FOREIGN KEY (book_id) REFERENCES books(id),
            FOREIGN KEY (genre_id) REFERENCES genres(id)
        );`,

		// Post genres table (Pivot)
		`CREATE TABLE IF NOT EXISTS post_genres (
            id INT AUTO_INCREMENT PRIMARY KEY,
            post_id INT NOT NULL,
            genre_id INT NOT NULL,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (genre_id) REFERENCES genres(id)
            );`,

		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			token VARCHAR(512) NOT NULL UNIQUE,
			device_info VARCHAR(255),  -- Optional: Device/browser details
			expires_at DATETIME NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS user_preferred_genres (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            genre_id INT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY (genre_id) REFERENCES genres(id)
        );`,
		// // Seller Reviews table
		// `CREATE TABLE IF NOT EXISTS seller_reviews (
		//     id INT AUTO_INCREMENT PRIMARY KEY,
		//     buyer_id INT NOT NULL,
		//     seller_id INT NOT NULL,
		//     rating INT NOT NULL,
		//     comment TEXT,
		//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		//     FOREIGN KEY (buyer_id) REFERENCES users(id),
		//     FOREIGN KEY (seller_id) REFERENCES users(id)
		// );`,
		// // Recommendations table
		// `CREATE TABLE IF NOT EXISTS recommendations (
		//     id INT AUTO_INCREMENT PRIMARY KEY,
		//     user_id INT NOT NULL,
		//     book_id INT NOT NULL,
		//     score DECIMAL(5,2) NOT NULL,
		//     generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		//     FOREIGN KEY (book_id) REFERENCES books(id)
		//     );`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error running migration query: %v", err)
		}
	}
	log.Println("Migrations executed successfully!")
}
