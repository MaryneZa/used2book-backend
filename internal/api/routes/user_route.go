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
)

func UserRoutes(db *sql.DB) http.Handler {


	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	uploadService := services.NewUploadService(userRepo)

	
	userHandler := &handlers.UserHandler{
		UserService:  userService,
		UploadService:  uploadService,
	}

	r := chi.NewRouter()

	r.With(middleware.AuthMiddleware).Get("/me", userHandler.GetMeHandler)
	
	r.With(middleware.AuthMiddleware).Post("/upload-profile-image", userHandler.UploadProfileImageHandler)
	r.With(middleware.AuthMiddleware).Post("/upload-background-image", userHandler.UploadBackgroundImageHandler)

	r.With(middleware.AuthMiddleware).Post("/edit-account-info", userHandler.EditAccountInfoHandler)
	r.With(middleware.AuthMiddleware).Post("/edit-username", userHandler.EditUserNameHandler)
	r.With(middleware.AuthMiddleware).Post("/edit-preferrence", userHandler.EditPreferrenceHandler)

	r.With(middleware.AuthMiddleware).Post("/add-library", userHandler.AddBookToLibraryHandler)

	r.With(middleware.AuthMiddleware).Get("/get-listing", userHandler.GetMyListingsHandler)
	r.With(middleware.AuthMiddleware).Get("/get-library", userHandler.GetMyLibraryHandler)
	r.With(middleware.AuthMiddleware).Get("/get-wishlist", userHandler.GetMyWishlist)

	r.With(middleware.AuthMiddleware).Get("/get-listing/{userID:[0-9]+}", userHandler.GetUserListingsHandler)
	r.With(middleware.AuthMiddleware).Get("/get-library/{userID:[0-9]+}", userHandler.GetUserLibraryHandler)
	r.With(middleware.AuthMiddleware).Get("/get-wishlist/{userID:[0-9]+}", userHandler.GetUserWishlist)


	r.Get("/user-info/{userID:[0-9]+}", userHandler.GetUserByIDHandler)
	r.With(middleware.AuthMiddleware).Get("/user-wishlist/{userID:[0-9]+}", userHandler.GetUserWishlist)


	r.With(middleware.AuthMiddleware).Get("/book-wishlist/{bookID:[0-9]+}", userHandler.AddBookToWishListHandler)
	r.With(middleware.AuthMiddleware).Get("/book-is-in-wishlist/{bookID:[0-9]+}", userHandler.IsBookInWishlistHandler)

	r.With(middleware.AuthMiddleware).Get("/get-listing-by-id/{listingID:[0-9]+}", userHandler.GetListingWithBookByIDHandler)

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

	r.With(middleware.AuthMiddleware).Post("/post-create", userHandler.CreatePostHandler)
	r.With(middleware.AuthMiddleware).Post("/upload-post-images", userHandler.UploadPostImagesHandler)

	r.With(middleware.AuthMiddleware).Get("/posts", userHandler.GetAllPostsHandler)
	r.With(middleware.AuthMiddleware).Get("/user-posts/{userID:[0-9]+}", userHandler.GetPostsByUserIDHandler)

	r.With(middleware.AuthMiddleware).Post("/user/comment-create", userHandler.CreateCommentHandler)
	r.With(middleware.AuthMiddleware).Get("/user/comments", userHandler.GetCommentsByPostIDHandler) // e.g., /user/comments?post_id=1

	// r.Get("/post", uh.GetPostByPostIDHandler) // e.g., /post?post_id=1

	r.Get("/user-review", userHandler.GetAllUserReview)






	
	
	return r

}
