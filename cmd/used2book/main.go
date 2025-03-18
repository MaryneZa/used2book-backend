package main

import (
	// "os"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"used2book-backend/internal/api"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/twiliootp" // adjust the import path to your module name and structure
	"used2book-backend/internal/utils"
	"used2book-backend/internal/services"

	
   
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	// log.Println("ENV - main" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - main")
    //     }
    // }

	

	

	db := utils.GetDB()

	utils.InitRedis()

	twiliootp.InitTwilio()

	// Assign the shared Redis client to the twiliootp package.
	twiliootp.RedisClient = utils.RedisClient

	// Initialize Twilio.

	// err := utils.InitFirebase()
	// if err != nil {
	//     log.Fatalf("cannot init firebase: %v", err)
	// }
    
	router := api.SetupRouter(db)



	utils.RunMigrations()

    // Initialize Book Repository & Service
    bookRepo := mysql.NewBookRepository(db)
    bookService := services.NewBookService(bookRepo)

    // ✅ Call SyncBooksIfNeeded() instead of inline logic
    bookService.SyncBooksIfNeeded()

	
	userRepo := mysql.NewUserRepository(db)

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background cleanup as a goroutine
	go userRepo.CleanupExpiredListings(ctx)



	log.Println("Server is listening on port 6951")
	if err := http.ListenAndServe(":6951", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

