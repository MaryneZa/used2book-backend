package handlers

import (
	"log"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "github.com/dchest/uniuri"
    "golang.org/x/oauth2"
    "used2book-backend/internal/config"
	"used2book-backend/internal/models"
	"used2book-backend/internal/services"
)




type AuthHandler struct{
	UserService *services.UserService
}

func (ah *AuthHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "<a href='/signup'>Sign up with Google</a>")
    fmt.Fprintln(w, "<a href='/login'>Log in with Google</a>")

}

// func (ah *AuthHandler) InitiateGoogleOAuth(w http.ResponseWriter, r *http.Request) {
//     oauthStateString := uniuri.New()
//     log.Printf("oauthStateString: %s", oauthStateString)
//     url := config.GoogleOauthConfig.AuthCodeURL(oauthStateString)
//     http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }


// func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request){
	
// 	InitiateGoogleOAuth(w, r)
	
// 	user, err := ah.UserService.Signup(r.Context(), reqUser)
	
// 	if err != nil {
// 		http.Error(w, "Signup failed: "+ err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Printf( "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

// 	fmt.Fprintf(w, "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(user)

// }

// func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

// 	InitiateGoogleOAuth(w, r)

// }

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
    oauthStateString := uniuri.New()
    log.Printf("Login oauthStateString: %s", oauthStateString)
    url := config.GoogleLoginConfig.AuthCodeURL(oauthStateString)
	log.Printf("Redirecting to: %s", url)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ah *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
    oauthStateString := uniuri.New()
    log.Printf("Signup oauthStateString: %s", oauthStateString)
    url := config.GoogleSignupConfig.AuthCodeURL(oauthStateString)
	log.Printf("Redirecting to: %s", url)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}


func (ah *AuthHandler) LoginCallbackHandler(w http.ResponseWriter, r *http.Request){
	code := r.FormValue("code")
    token, err := config.GoogleLoginConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        http.Error(w, "Code exchange failed", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Token: %s\n", token.AccessToken)

    response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    contents, err := io.ReadAll(response.Body)
    if err != nil {
        http.Error(w, "Failed to read user info", http.StatusInternalServerError)
        return
    }
    // var user models.User

	var user models.LoginUser

    if err := json.Unmarshal(contents, &user); err != nil {
        http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
        return
    }

	userEmail, err := ah.UserService.Login(r.Context(), user)

	if err != nil {
		http.Error(w, "Signup failed: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf( "Email: %s", userEmail)

    fmt.Fprintf(w, "login successfully!")

    fmt.Fprintf(w, "Email: %s", userEmail)

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)

}

func (ah *AuthHandler) SignupCallbackHandler(w http.ResponseWriter, r *http.Request){
	code := r.FormValue("code")
    token, err := config.GoogleSignupConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        http.Error(w, "Code exchange failed", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "Token: %s\n", token.AccessToken)

    response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    contents, err := io.ReadAll(response.Body)
    if err != nil {
        http.Error(w, "Failed to read user info", http.StatusInternalServerError)
        return
    }
    // var user models.User

	var reqUser models.SignupUser

    if err := json.Unmarshal(contents, &reqUser); err != nil {
        http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
        return
    }

	user, err := ah.UserService.Signup(r.Context(), reqUser)

	if err != nil {
		http.Error(w, "Signup failed: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf( "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

	fmt.Fprintf(w, "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// func (ah *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request){
// 	code := r.FormValue("code")
//     token, err := config.GoogleOauthConfig.Exchange(oauth2.NoContext, code)
//     if err != nil {
//         http.Error(w, "Code exchange failed", http.StatusInternalServerError)
//         return
//     }
//     fmt.Fprintf(w, "Token: %s\n", token.AccessToken)

//     response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
//     if err != nil {
//         http.Error(w, "Failed to get user info", http.StatusInternalServerError)
//         return
//     }
//     defer response.Body.Close()

//     contents, err := io.ReadAll(response.Body)
//     if err != nil {
//         http.Error(w, "Failed to read user info", http.StatusInternalServerError)
//         return
//     }
//     // var user models.User

// 	var user models.OathUser

//     if err := json.Unmarshal(contents, &user); err != nil {
//         http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
//         return
//     }

// 	log.Printf( "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

//     fmt.Fprintf(w, "Email: %s\nName: %s\nImage link: %s\n", user.Email, user.Name, user.ProfilePictureURL)

	
// }
