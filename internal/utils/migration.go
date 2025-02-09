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

			username VARCHAR(255) DEFAULT NULL,
			provider ENUM('google','local') NOT NULL,
            
            hashed_password VARCHAR(255),
            phone_number VARCHAR(20) DEFAULT NULL,
			picture VARCHAR(255) DEFAULT NULL,
			role ENUM('user','admin') DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
            );`,
            
        // provider_id VARCHAR(255),
        

        // Books table
        `CREATE TABLE IF NOT EXISTS books (
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255),
            description TEXT,
            category VARCHAR(255),
            cover_image_url VARCHAR(255),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );`,

        // User Libraries table
        `CREATE TABLE IF NOT EXISTS user_libraries (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            book_id INT NOT NULL,
            status ENUM('own','read','wishlist') NOT NULL,
            personal_notes TEXT,
            favorite_quotes TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

        // Listings table
        `CREATE TABLE IF NOT EXISTS listings (
            id INT AUTO_INCREMENT PRIMARY KEY,
            seller_id INT NOT NULL,
            book_id INT NOT NULL,
            price DECIMAL(10,2) NOT NULL,
            status ENUM('for_sale','sold','removed') DEFAULT 'for_sale',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (seller_id) REFERENCES users(id),
            FOREIGN KEY (book_id) REFERENCES books(id)
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
            FOREIGN KEY (buyer_id) REFERENCES users(id)
        );`,

        // Transactions table
        `CREATE TABLE IF NOT EXISTS transactions (
            id INT AUTO_INCREMENT PRIMARY KEY,
            buyer_id INT NOT NULL,
            seller_id INT NOT NULL,
            listing_id INT NOT NULL,
            transaction_amount DECIMAL(10,2) NOT NULL,
            payment_status ENUM('pending','completed','failed') DEFAULT 'pending',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (buyer_id) REFERENCES users(id),
            FOREIGN KEY (seller_id) REFERENCES users(id),
            FOREIGN KEY (listing_id) REFERENCES listings(id)
        );`,

        // Book Reviews table
        `CREATE TABLE IF NOT EXISTS book_reviews (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            book_id INT NOT NULL,
            rating INT NOT NULL,
            comment TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

        // Seller Reviews table
        `CREATE TABLE IF NOT EXISTS seller_reviews (
            id INT AUTO_INCREMENT PRIMARY KEY,
            buyer_id INT NOT NULL,
            seller_id INT NOT NULL,
            rating INT NOT NULL,
            comment TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (buyer_id) REFERENCES users(id),
            FOREIGN KEY (seller_id) REFERENCES users(id)
        );`,

        // Posts table
        `CREATE TABLE IF NOT EXISTS posts (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            title VARCHAR(255) NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
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
            FOREIGN KEY (user_id) REFERENCES users(id)
        );`,

        // Notifications table
        `CREATE TABLE IF NOT EXISTS notifications (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            message VARCHAR(255),
            is_read BOOLEAN DEFAULT false,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );`,

        // Recommendations table
        `CREATE TABLE IF NOT EXISTS recommendations (
            id INT AUTO_INCREMENT PRIMARY KEY,
            user_id INT NOT NULL,
            book_id INT NOT NULL,
            score DECIMAL(5,2) NOT NULL,
            generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (book_id) REFERENCES books(id)
        );`,

        // Tags table
        `CREATE TABLE IF NOT EXISTS tags (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        );`,

        // Book Tags table (Pivot)
        `CREATE TABLE IF NOT EXISTS book_tags (
            id INT AUTO_INCREMENT PRIMARY KEY,
            book_id INT NOT NULL,
            tag_id INT NOT NULL,
            FOREIGN KEY (book_id) REFERENCES books(id),
            FOREIGN KEY (tag_id) REFERENCES tags(id)
        );`,

        // Post Tags table (Pivot)
        `CREATE TABLE IF NOT EXISTS post_tags (
            id INT AUTO_INCREMENT PRIMARY KEY,
            post_id INT NOT NULL,
            tag_id INT NOT NULL,
            FOREIGN KEY (post_id) REFERENCES posts(id),
            FOREIGN KEY (tag_id) REFERENCES tags(id)
        );`,
        `CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			token VARCHAR(512) NOT NULL UNIQUE,
			device_info VARCHAR(255),  -- Optional: Device/browser details
			expires_at DATETIME NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`,
    }

    for _, query := range queries {
        _, err := db.Exec(query)
        if err != nil {
            log.Fatalf("Error running migration query: %v", err)
        }
    }
    log.Println("Migrations executed successfully!")
}
