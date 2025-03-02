package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"used2book-backend/internal/services"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	OmiseService *services.OmiseService
	UserService  *services.UserService
}

// ChargeRequest represents the JSON request body structure
type ChargeRequest struct {
	ListingID int    `json:"listing_id"`
	Token     string `json:"token"`
	BuyerID   int    `json:"buyer_id"`
}

// ChargeHandler processes a payment from the buyer to the seller
func (ph *PaymentHandler) ChargeHandler(w http.ResponseWriter, r *http.Request) {
	var req ChargeRequest

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// 1️⃣ Fetch listing details (Get seller info)
	listing, err := ph.UserService.GetListingByID(context.Background(), req.ListingID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Listing not found")
		return
	}

	// 2️⃣ Ensure the seller has a valid Omise account
	if !listing.SellerOmiseID.Valid {
		sendErrorResponse(w, http.StatusBadRequest, "Seller does not have a valid Omise account")
		return
	}
	sellerOmiseID := listing.SellerOmiseID.String // Extract the string value

	// 3️⃣ Process payment using Omise
	charge, err := ph.OmiseService.CreateCharge(int64(listing.Price*100), req.Token, sellerOmiseID)
	if err != nil {
		log.Println("❌ Omise Charge Error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Payment processing failed")
		return
	}

	// 4️⃣ Mark the listing as "sold" after successful charge
	err = ph.UserService.MarkListingAsSold(context.Background(), req.ListingID, req.BuyerID, listing.Price)
	if err != nil {
		log.Println("❌ MarkListingAsSold Error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update listing status")
		return
	}

	// 5️⃣ Return success response
	response := map[string]interface{}{
		"success": true,
		"charge":  charge,
		"message": "Payment successful and listing marked as sold!",
	}
	sendSuccessResponse(w, response)
}


func (ph *PaymentHandler) CreateOmiseAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (assuming authentication middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode the request body
	var req struct {
		BankAccountNumber string `json:"bank_account_number"`
		BankAccountName   string `json:"bank_account_name"`
		BankCode          string `json:"bank_code"` // Example: "bbl" (Bangkok Bank), "scb" (Siam Commercial Bank)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println(req)

	// Ensure required fields are present
	if req.BankAccountNumber == "" || req.BankAccountName == "" || req.BankCode == "" {
		http.Error(w, "Bank details are required", http.StatusBadRequest)
		return
	}

	// Get user from the database
	user, err := ph.UserService.GetMe(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Ensure the user does not already have an Omise account
	if user.OmiseAccountID.Valid {
		http.Error(w, "User already has an Omise account", http.StatusConflict)
		return
	}

	// Create Omise recipient (bank account registration)
	recipientID, err := ph.OmiseService.CreateRecipient(req.BankAccountNumber, req.BankAccountName, req.BankCode)
	if err != nil {
		log.Println("Error creating Omise recipient:", err)
		http.Error(w, "Failed to create Omise account", http.StatusInternalServerError)
		return
	}

	// Store the Omise account ID in the database
	err = ph.UserService.UpdateOmiseAccountID(r.Context(), userID, recipientID)
	if err != nil {
		log.Println("Error saving Omise account ID to DB:", err)
		http.Error(w, "Failed to store Omise account ID", http.StatusInternalServerError)
		return
	}

	// Send success response
	resp := map[string]interface{}{
		"success":          true,
		"omise_recipient": recipientID,
	}
	json.NewEncoder(w).Encode(resp)
}
