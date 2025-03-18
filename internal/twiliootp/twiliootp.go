package twiliootp

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
    "errors"
	
	"github.com/go-redis/redis/v8"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
    "github.com/joho/godotenv"

)



func getTwilioAccountSID() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }
	// log.Println("ENV - twiloID" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - twilioID")
    //     }
    // }
    secret := os.Getenv("TWILIO_ACCOUNT_SID")
    
    if secret == "" {
        return "", errors.New("TWILIO_ACCOUNT_SID is not set in .env file")
    }
    return secret, nil
}

func getTwilioAuthToken() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }
	// log.Println("ENV - twilio token" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - twiliotoken")
    //     }
    // }
    secret := os.Getenv("TWILIO_AUTH_TOKEN")
    
    if secret == "" {
        return "", errors.New("TWILIO_AUTH_TOKEN is not set in .env file")
    }
    return secret, nil
}

func getTwilioPhoneNumber() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }
	// log.Println("ENV - twilio phone" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - twilio phone")
    //     }
    // }
    secret := os.Getenv("TWILIO_PHONE_NUMBER")
    
    if secret == "" {
        return "", errors.New("TWILIO_PHONE_NUMBER is not set in .env file")
    }
    return secret, nil
}

// twilioClient is the initialized Twilio REST client.
var twilioClient *twilio.RestClient

// RedisClient is the Redis client (assumed to be initialized elsewhere).
var RedisClient *redis.Client

// InitTwilio initializes the Twilio REST client.
func InitTwilio() {
	acc_phone, err := getTwilioPhoneNumber()
	if err != nil {
        fmt.Errorf("failed to get TWILIO_PHONE_NUMBER: %v", err)
    } 
	fmt.Println("Twilio Phone Number:", acc_phone)
	username, err := getTwilioAccountSID()
	if err != nil {
        fmt.Errorf("failed to get TWILIO_ACCOUNT_SID: %v", err)
    } 
	auth_token, err := getTwilioAuthToken()
	if err != nil {
        fmt.Errorf("failed to get TWILIO_AUTH_TOKEN: %v", err)
    } 
	
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: auth_token,
	})
}

// SendOTP generates a 6‑digit OTP, stores it in Redis for 5 minutes, and sends it via Twilio SMS.
func SendOTP(ctx context.Context, phoneNumber string) error {
	// Generate a random 6‑digit OTP.
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store the OTP in Redis with a 5‑minute expiration.
	if err := RedisClient.Set(ctx, "otp:"+phoneNumber, otp, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}

	acc_phone, err := getTwilioPhoneNumber()
	if err != nil {
        fmt.Errorf("failed to get TWILIO_PHONE_NUMBER: %v", err)
    } 

	// Prepare the SMS message.
	messageBody := fmt.Sprintf("Your OTP code is: %s", otp)
	params := &openapi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(acc_phone)
	params.SetBody(messageBody)

	// Send the SMS using Twilio.
	resp, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}
	fmt.Printf("SMS sent successfully. SID: %s\n", *resp.Sid)
	return nil
}

// VerifyOTP checks the OTP provided by the user against the stored value in Redis.
func VerifyOTP(ctx context.Context, phoneNumber, userOTP string) (bool, error) {
	storedOTP, err := RedisClient.Get(ctx, "otp:"+phoneNumber).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("OTP not found or expired")
	} else if err != nil {
		return false, fmt.Errorf("error retrieving OTP from Redis: %w", err)
	}

	if storedOTP == userOTP {
		// Optionally delete the OTP after successful verification.
		RedisClient.Del(ctx, "otp:"+phoneNumber)
		return true, nil
	}

	return false, nil
}

// ResendOTP simply calls SendOTP without any cooldown logic.
func ResendOTP(ctx context.Context, phoneNumber string) error {
	return SendOTP(ctx, phoneNumber)
}
