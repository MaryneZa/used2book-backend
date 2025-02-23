package utils

import (
    "context"
    "fmt"
    "log"

    firebase "firebase.google.com/go"
    "google.golang.org/api/option"

    // 1) Import the embed package
    _ "embed"
)

// 2) Embed the file into a byte slice at compile time
//go:embed used2book-otp-firebase-secretKey.json
var firebaseCredentials []byte

// FirebaseApp will hold the initialized Firebase app
var FirebaseApp *firebase.App

// InitFirebase initializes Firebase using the embedded credentials.
func InitFirebase() error {
    // 3) Pass the embedded JSON bytes to WithCredentialsJSON
    app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(firebaseCredentials))
    if err != nil {
        log.Fatalf("Failed to initialize Firebase: %v", err)
        return fmt.Errorf("failed to initialize Firebase app: %w", err)
    }

    FirebaseApp = app
    fmt.Println("🔥 Firebase Initialized Successfully (using go:embed)")
    return nil
}
