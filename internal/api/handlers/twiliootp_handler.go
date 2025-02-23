package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"used2book-backend/internal/services"
	"used2book-backend/internal/twiliootp" // adjust the import path to your module name and structure
)

type TwilioOTPHandler struct {
	UserService *services.UserService
}

// OTPRequest represents the JSON payload for sending an OTP.
type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// OTPVerifyRequest represents the JSON payload for verifying an OTP.
type OTPVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

func (th *TwilioOTPHandler) SendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req OTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// ✅ Check if UserService is nil
	if th.UserService == nil {
		log.Println("❌ Error: UserService is nil!")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ✅ Log incoming request
	log.Println("📞 Request to send OTP for phone:", req.PhoneNumber)

	// ✅ Check if the phone number exists
	taken, err := th.UserService.IsPhoneNumberTaken(ctx, req.PhoneNumber)
	if err != nil {
		log.Println("❌ Error checking phone number:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "phone number internal server error: "+err.Error())
		return
	}

	// ✅ Ensure `taken` is correctly handled
	if taken {
		log.Println("❌ Phone number already registered:", req.PhoneNumber)
		sendErrorResponse(w, http.StatusConflict, "phone number already registered")
		return
	}


	// ✅ Send OTP and log response
	err = twiliootp.SendOTP(ctx, req.PhoneNumber)
	if err != nil {
		log.Println("❌ Error sending OTP:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to send OTP: "+err.Error())
		return
	}

	log.Println("✅ OTP sent successfully to:", req.PhoneNumber)

	// ✅ Success response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "OTP sent successfully",
	})
}


// verifyOTPHandler handles the route for verifying an OTP.
func (th *TwilioOTPHandler) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req OTPVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	valid, err := twiliootp.VerifyOTP(ctx, req.PhoneNumber, req.OTP)
	if err != nil {
		log.Printf("Error verifying OTP: %v", err)
		http.Error(w, fmt.Sprintf("Failed to verify OTP: %v", err), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("user_id").(int)

	// 2. Check if user with same email already exists
	err = th.UserService.EditPhoneNumber(r.Context(), userID, req.PhoneNumber)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "phone number "+err.Error()) // 409 Conflict if user exists
		return
	}

	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"success": valid,
		"message": func() string {
			if valid {
				return "OTP verified successfully"
			}
			return "OTP verification failed"
		}(),
	})

}

// (Optional) If you want a separate resend endpoint that simply calls SendOTP:
func (th *TwilioOTPHandler) ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	// You might use the same payload as sending an OTP.
	th.SendOTPHandler(w, r)
}
