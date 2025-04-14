package routes

import (
	"database/sql"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"
	"used2book-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver

	"github.com/streadway/amqp"

)


func UserRoutes(db *sql.DB, rabbitConn *amqp.Connection) http.Handler {


	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	uploadService := services.NewUploadService(userRepo)

	
	userHandler := &handlers.UserHandler{
		UserService:  userService,
		UploadService:  uploadService,
		RabbitMQConn: rabbitConn,
	}

	r := chi.NewRouter()

	r.Get("/test", userHandler.TestPort)


	r.With(middleware.AuthMiddleware).Get("/me", userHandler.GetMeHandler)

	r.With(middleware.AuthMiddleware).Post("/create-bank-account", userHandler.CreateBankAccountHandler)

	
	r.With(middleware.AuthMiddleware).Post("/upload-profile-image", userHandler.UploadProfileImageHandler)
	r.With(middleware.AuthMiddleware).Post("/upload-background-image", userHandler.UploadBackgroundImageHandler)

	r.With(middleware.AuthMiddleware).Post("/edit-account-info", userHandler.EditAccountInfoHandler)
	r.With(middleware.AuthMiddleware).Post("/edit-username", userHandler.EditUserNameHandler)
	r.With(middleware.AuthMiddleware).Post("/edit-profile", userHandler.EditProfileHandler)

	r.With(middleware.AuthMiddleware).Post("/add-library", userHandler.AddBookToLibraryHandler)

	r.With(middleware.AuthMiddleware).Post("/add-listing", userHandler.AddBookToListingHandler)


	r.With(middleware.AuthMiddleware).Get("/get-listing", userHandler.GetMyListingsHandler)
	r.With(middleware.AuthMiddleware).Get("/get-library", userHandler.GetMyLibraryHandler)
	r.With(middleware.AuthMiddleware).Get("/get-wishlist", userHandler.GetMyWishlist)

	r.With(middleware.AuthMiddleware).Get("/get-listing/{userID:[0-9]+}", userHandler.GetUserListingsHandler)
	r.With(middleware.AuthMiddleware).Get("/get-library/{userID:[0-9]+}", userHandler.GetUserLibraryHandler)
	r.Get("/get-wishlist/{userID:[0-9]+}", userHandler.GetUserWishlist)

	r.With(middleware.AuthMiddleware).Post("/listing/remove/{listingID:[0-9]+}", userHandler.RemoveListingHandler)

	r.Get("/user-info/{userID:[0-9]+}", userHandler.GetUserByIDHandler)


	r.With(middleware.AuthMiddleware).Get("/book-wishlist/{bookID:[0-9]+}", userHandler.AddBookToWishListHandler)
	r.With(middleware.AuthMiddleware).Get("/book-is-in-wishlist/{bookID:[0-9]+}", userHandler.IsBookInWishlistHandler)

	r.With(middleware.AuthMiddleware).Get("/get-listing-by-id/{listingID:[0-9]+}", userHandler.GetListingByIDHandler)

	r.Get("/all-users", userHandler.GetAllUsersHandler)

	r.With(middleware.AuthMiddleware).With(middleware.AdminMiddleware(db)).Get("/user-count", userHandler.GetUserCount) // Sync books from Google Sheets

	// r.With(middleware.AuthMiddleware).Post("/edit-phone-number", userHandler.EditPhoneNumberHandler)

	r.With(middleware.AuthMiddleware).Post("/listing/sold", userHandler.MarkListingAsSoldHandler)

	r.With(middleware.AuthMiddleware).Post("/preferences", userHandler.SetUserPreferredGenresHandler)
	r.With(middleware.AuthMiddleware).Get("/preferences", userHandler.GetUserPreferencesHandler)

	r.Get("/user-preferences", userHandler.GetAllUserPreferred)


	r.With(middleware.AuthMiddleware).Post("/gender", userHandler.UpdateGenderHandler)
	r.With(middleware.AuthMiddleware).Get("/gender", userHandler.GetGenderHandler)

	r.With(middleware.AuthMiddleware).Post("/cart", userHandler.AddToCartHandler)
	r.With(middleware.AuthMiddleware).Get("/cart", userHandler.GetCartHandler)
	r.With(middleware.AuthMiddleware).Post("/cart-rm", userHandler.RemoveFromCartHandler)

	r.With(middleware.AuthMiddleware).Post("/offers", userHandler.AddToOffersHandler)
	r.With(middleware.AuthMiddleware).Get("/offers/buyer", userHandler.GetBuyerOffersHandler)
	r.With(middleware.AuthMiddleware).Get("/offers/seller", userHandler.GetSellerOffersHandler)
    r.With(middleware.AuthMiddleware).Post("/offers-rm", userHandler.RemoveFromOffersHandler)
	r.With(middleware.AuthMiddleware).Post("/offers/accept", userHandler.AcceptOfferHandler)
	r.With(middleware.AuthMiddleware).Post("/offers/reject", userHandler.RejectOfferHandler)
	r.With(middleware.AuthMiddleware).Get("/offers/{offerID:[0-9]+}/payment", userHandler.GetAcceptedOfferHandler)

	r.With(middleware.AuthMiddleware).Post("/post-create", userHandler.CreatePostHandler)
	r.With(middleware.AuthMiddleware).Post("/upload-post-images", userHandler.UploadPostImagesHandler)

	r.With(middleware.AuthMiddleware).Get("/posts", userHandler.GetAllPostsHandler)
	r.With(middleware.AuthMiddleware).Get("/user-posts/{userID:[0-9]+}", userHandler.GetPostsByUserIDHandler)

	r.With(middleware.AuthMiddleware).Post("/comment-create", userHandler.CreateCommentHandler)
	r.With(middleware.AuthMiddleware).Get("/comments/{postID:[0-9]+}", userHandler.GetCommentsByPostIDHandler) 

	r.With(middleware.AuthMiddleware).Post("/like-toggle/{postID:[0-9]+}", userHandler.ToggleLikeHandler)
	r.With(middleware.AuthMiddleware).Get("/like-count/{postID:[0-9]+}", userHandler.GetLikeCountHandler) 
	r.With(middleware.AuthMiddleware).Get("/like-check/{postID:[0-9]+}", userHandler.IsPostLikedHandler) 

	r.With(middleware.AuthMiddleware).Get("/all", userHandler.GetAllUsersHandler) 

	r.With(middleware.AuthMiddleware).Get("/purchase-listing", userHandler.GetMyPurchasedListingsHandler) 

	r.With(middleware.AuthMiddleware).Get("/my-orders", userHandler.GetMyOrdersHandler)
	
	r.With(middleware.AuthMiddleware).Get("/user-wishlist/{bookID:[0-9]+}", userHandler.GetUsersByWishlistBookIDHandler) 
	
	r.With(middleware.AuthMiddleware).Post("/delete-user-library/{id:[0-9]+}", userHandler.DeleteUserLibraryByIDHandler) 

	r.With(middleware.AuthMiddleware).Post("/book-request", userHandler.CreateBookRequestHandle) 

	r.With(middleware.AuthMiddleware).Get("/book-request", userHandler.GetBookRequestHandler) 

	// r.Get("/post", uh.GetPostByPostIDHandler) // e.g., /post?post_id=1

	r.Get("/user-review", userHandler.GetAllUserReview)

	// main.go or router.go






	
	
	return r

}
