package handlers

import (
    "time"
    "net/http"
	"used2book-backend/internal/utils"
	"used2book-backend/internal/services"
    "strconv"
)


type TokenHandler struct{
	TokenService *services.TokenService
	UserService *services.UserService
}

func (th *TokenHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
    // Step 1: Get the refresh token from the cookie
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        sendErrorResponse(w, http.StatusUnauthorized, "Refresh token missing")
        return
    }
    refreshToken := cookie.Value

	userID, err := th.TokenService.ValidateRefreshToken(r.Context(), refreshToken)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Authentication failed(refresh_token): "+err.Error()) // 409 Conflict if user exists
		return
	}
    // Convert userID to string
    userIDStr := strconv.Itoa(userID)

	accessToken, refreshToken, err := th.TokenService.GenerateTokens(r.Context(), userIDStr)
	if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Error generating new token")
        return
    }

    err = th.TokenService.UpdateRefreshToken(r.Context(), userID, refreshToken)
	if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Error updating new token")
        return
    }

    // Step 7: Set cookies with new tokens
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Expires:  time.Now().Add(15 * time.Minute),
    })

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Expires:  utils.RefreshTokenExpiration(),
    })

    // Step 8: Send response with new tokens
    sendSuccessResponse(w, map[string]interface{}{
        "success":      true,
        "access_token": accessToken,
        "refresh_token": refreshToken,
    })
}

