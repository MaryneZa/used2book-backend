// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"todo-api/internal/routes"
// 	"todo-api/internal/utils"
//     "github.com/go-chi/chi/v5"
// 	"log"
// 	"github.com/joho/godotenv"
// )

// func main() {

//     utils.GetDB()
// 	utils.RunMigrations()

// 	fmt.Println("Server is running on port 8080...")

//     // Start the server
// 	if err := http.ListenAndServe(":8080", r); err != nil {
// 		log.Fatalf("Error starting server: %v", err)
// 	}

// 	// data.TestConnection()
// }
package main

import (
    "log"
    "net/http"
    "used2book-backend/internal/api"
    "used2book-backend/internal/utils"
	"github.com/joho/godotenv"
)



func main() {

	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    db := utils.GetDB()

    router := api.SetupRouter(db)

	utils.RunMigrations()

    log.Println("Server is listening on port 6951")
    if err := http.ListenAndServe(":6951", router); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
