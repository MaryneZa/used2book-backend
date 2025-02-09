package utils

import (
    "errors"
    "github.com/golang-jwt/jwt/v4"
    "github.com/joho/godotenv"
    "os"
    "time"
)

// Load the secret key from .env
func getAccessTokenSecretKey() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }
    secret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
    
    if secret == "" {
        return "", errors.New("JWT_SECRET is not set in .env file")
    }
    return secret, nil
}

func getRefreshTokenSecretKey() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }
    secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
    if secret == "" {
        return "", errors.New("JWT_SECRET is not set in .env file")
    }
    return secret, nil
}





// GenerateAccessTokenJWT generates a JWT token for a given user ID
func GenerateAccessToken(userID int) (string, error) {
    secretKey, err := getAccessTokenSecretKey()
    if err != nil {
        return "", err
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "type":    "access",
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    return token.SignedString([]byte(secretKey))
}

// Generate Refresh Token
func GenerateRefreshToken(userID int) (string, error) {
    secretKey, err := getRefreshTokenSecretKey()
    if err != nil {
        return "", err
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "type":    "refresh",
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // Expires in 7 days
    })

    return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenStr string, tokenType string) (int, error) {
    var secretKey string
    var err error

    if tokenType == "access"{
        secretKey, err = getAccessTokenSecretKey()
    }else if tokenType == "refresh"{
        secretKey, err = getRefreshTokenSecretKey()
    }else{
        return 0, errors.New("invalid token type")
    }

    if err != nil {
        return 0, err
    }

    // Parse and validate the token
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return 0, errors.New("token has expired")
        }
        return 0, err
    }

    // Extract claims
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if claims["type"] != tokenType { // Ensure it matches the expected type
            return 0, errors.New("invalid token type")
        }

        userID, ok := claims["user_id"].(float64)
        if !ok {
            return 0, errors.New("invalid user_id claim")
        }
        return int(userID), nil
    }

    return 0, errors.New("invalid token")
}

func RefreshTokenExpiration() time.Time {
    return time.Now().Add(7 * 24 * time.Hour) // 7 days expiry
}
