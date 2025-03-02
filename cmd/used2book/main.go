package main

import (
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



	log.Println("Server is listening on port 6951")
	if err := http.ListenAndServe(":6951", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

