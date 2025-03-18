package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"used2book-backend/internal/services"

	"github.com/omise/omise-go"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	OmiseService *services.OmiseService
	UserService  *services.UserService
}

// ChargeRequest represents the JSON request body structure
type ChargeRequest struct {
	ListingID int `json:"listing_id"`
	BuyerID   int `json:"buyer_id"`
}

// ChargeHandler processes a PromptPay payment from the buyer to the seller
func (ph *PaymentHandler) ChargeHandler(w http.ResponseWriter, r *http.Request) {
	var req ChargeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	log.Println(req.BuyerID, req.ListingID)
	// Fetch listing details upfront
	listing, err := ph.UserService.GetListingByID(context.Background(), req.ListingID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "Listing not found")
		return
	}

	// Attempt to reserve the listing atomically
	success, err := ph.UserService.ReserveListing(context.Background(), req.ListingID, 1)
	if err != nil {
		log.Println("❌ ReserveListing Error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to process purchase")
		return
	}
	if !success {
		// Check if the listing is reserved and its expiration status
		isReserved, isExpired, err := ph.UserService.IsListingReserved(context.Background(), req.ListingID)
		if err != nil {
			log.Println("❌ Error checking reservation status:", err)
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to check listing status")
			return
		}

		if isReserved && !isExpired {
			// Active reservation
			sendErrorResponse(w, http.StatusConflict, "This book is currently reserved by another buyer. Try again later.")
			return
		} else if isReserved && isExpired {
			// Expired reservation, attempt to expire and retry
			if err := ph.UserService.ExpireReservedListing(context.Background(), req.ListingID); err != nil {
				log.Println("❌ ExpireReservedListing Error:", err)
				sendErrorResponse(w, http.StatusInternalServerError, "Failed to process listing status")
				return
			}
			success, err = ph.UserService.ReserveListing(context.Background(), req.ListingID, 1)
			if err != nil || !success {
				sendErrorResponse(w, http.StatusConflict, "Book is no longer available")
				return
			}
		} else {
			// Listing exists but isn’t for_sale or reserved_sale (e.g., sold, removed)
			sendErrorResponse(w, http.StatusConflict, "Book is not available for sale")
			return
		}
	}

	// Verify seller’s Omise account in the listing
	if !listing.SellerOmiseID.Valid {
		_ = ph.UserService.ExpireReservedListing(context.Background(), req.ListingID)
		sendErrorResponse(w, http.StatusBadRequest, "Seller does not have a valid Omise account for this listing")
		return
	}
	sellerOmiseID := listing.SellerOmiseID.String

	// Create PromptPay charge with 1-minute expiration
	charge, err := ph.OmiseService.CreatePromptPayCharge(int64(listing.Price*100), req.ListingID, sellerOmiseID, req.BuyerID, 2)
	if err != nil {
		_ = ph.UserService.ExpireReservedListing(context.Background(), req.ListingID)
		log.Println("❌ PromptPay Charge Error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Payment initiation failed")
		return
	}

	// Record transaction
	err = ph.UserService.CreateTransaction(context.Background(), req.BuyerID, listing.SellerID, req.ListingID, float64(listing.Price), "reserved")
	if err != nil {
		_ = ph.UserService.ExpireReservedListing(context.Background(), req.ListingID)
		if strings.Contains(err.Error(), "Duplicate entry") {
			sendErrorResponse(w, http.StatusConflict, "Transaction already initiated for this book")
			return
		}
		log.Println("❌ CreateTransaction Error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to record transaction")
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"charge_id":  charge.ID,
		"qr_code":    charge.Source.ScannableCode.Image.DownloadURI,
		"amount":     float64(charge.Amount) / 100,
		"expires_at": charge.ExpiresAt,
		"message":    "Scan this QR code within 1 minutes to pay. The book is reserved for you.",
	}
	sendSuccessResponse(w, response)
}

// WebhookHandler processes Omise payment confirmations
func (ph *PaymentHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var event struct {
		Key  string         `json:"key"`
		Data *omise.Charge `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println("❌ Webhook Decode Error:", err)
		http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
		return
	}

	if event.Key == "charge.complete" && event.Data.Paid {
		listingID, ok := event.Data.Metadata["listing_id"].(float64)
		if !ok {
			log.Println("❌ Invalid listing_id in metadata")
			http.Error(w, "Invalid metadata", http.StatusBadRequest)
			return
		}
		buyerID, ok := event.Data.Metadata["buyer_id"].(float64)
		if !ok {
			log.Println("❌ BuyerID not found in metadata")
			http.Error(w, "Invalid metadata", http.StatusBadRequest)
			return
		}
		amount := float64(event.Data.Amount / 100)

		// Update transaction status
		err := ph.UserService.UpdateTransactionStatus(context.Background(), int(listingID), "completed")
		if err != nil {
			log.Println("❌ UpdateTransactionStatus Error:", err)
			http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
			return
		}

		// Mark listing as sold
		err = ph.UserService.MarkListingAsSold(context.Background(), int(listingID), int(buyerID), amount)
		if err != nil {
			log.Println("❌ MarkListingAsSold Error:", err)
			http.Error(w, "Failed to update listing", http.StatusInternalServerError)
			return
		}

		log.Printf("Payment confirmed for listing %d by buyer %d", int(listingID), int(buyerID))
	} else if event.Key == "charge.expired" {
		listingID, ok := event.Data.Metadata["listing_id"].(float64)
		if !ok {
			log.Println("❌ Invalid listing_id in metadata")
			http.Error(w, "Invalid metadata", http.StatusBadRequest)
			return
		}
		// Expire the listing and transaction
		err := ph.UserService.ExpireReservedListing(context.Background(), int(listingID))
		if err != nil {
			log.Println("❌ ExpireReservedListing Error:", err)
			http.Error(w, "Failed to expire listing", http.StatusInternalServerError)
			return
		}
		log.Printf("Charge expired for listing %d", int(listingID))
	}

	w.WriteHeader(http.StatusOK)
}

// CreateOrUpdateOmiseAccountHandler sets up or updates a seller's Omise recipient
func (ph *PaymentHandler) CreateOrUpdateOmiseAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		BankAccountNumber string `json:"bank_account_number"`
		BankAccountName   string `json:"bank_account_name"`
		BankCode          string `json:"bank_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.BankAccountNumber == "" || req.BankAccountName == "" || req.BankCode == "" {
		http.Error(w, "Bank details are required", http.StatusBadRequest)
		return
	}

	user, err := ph.UserService.GetMe(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var recipientID string
	if user.OmiseAccountID.Valid {
		_, err := ph.OmiseService.GetRecipient(user.OmiseAccountID.String)
		if err != nil {
			if strings.Contains(err.Error(), "not_found") {
				log.Printf("Recipient %s not found, creating a new one", user.OmiseAccountID.String)
				recipientID, err = ph.OmiseService.CreateRecipient(req.BankAccountNumber, req.BankAccountName, req.BankCode)
				if err != nil {
					log.Println("Error creating Omise recipient:", err)
					http.Error(w, "Failed to create Omise account", http.StatusInternalServerError)
					return
				}
			} else {
				log.Println("Error retrieving Omise recipient:", err)
				http.Error(w, "Failed to verify Omise account", http.StatusInternalServerError)
				return
			}
		} else {
			recipientID, err = ph.OmiseService.UpdateRecipient(user.OmiseAccountID.String, req.BankAccountNumber, req.BankAccountName, req.BankCode)
			if err != nil {
				log.Println("Error updating Omise recipient:", err)
				http.Error(w, "Failed to update Omise account", http.StatusInternalServerError)
				return
			}
		}
	} else {
		recipientID, err = ph.OmiseService.CreateRecipient(req.BankAccountNumber, req.BankAccountName, req.BankCode)
		if err != nil {
			log.Println("Error creating Omise recipient:", err)
			http.Error(w, "Failed to create Omise account", http.StatusInternalServerError)
			return
		}
	}

	err = ph.UserService.UpdateOmiseAccountID(r.Context(), userID, recipientID)
	if err != nil {
		log.Println("Error saving Omise account ID to DB:", err)
		http.Error(w, "Failed to store Omise account ID", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"success":         true,
		"omise_recipient": recipientID,
	}
	json.NewEncoder(w).Encode(resp)
}

// GetBankAccountInfoHandler retrieves the user's bank account details from Omise
func (ph *PaymentHandler) GetBankAccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := ph.UserService.GetMe(r.Context(), userID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if !user.OmiseAccountID.Valid {
		sendErrorResponse(w, http.StatusBadRequest, "No bank account linked to this user")
		return
	}

	recipient, err := ph.OmiseService.GetRecipient(user.OmiseAccountID.String)
	if err != nil {
		log.Println("❌ Error retrieving Omise recipient:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve bank account info")
		return
	}

	bankAccount := recipient.BankAccount
	response := map[string]interface{}{
		"success":           true,
		"bank_account_number": bankAccount.Number,
		"bank_account_name":   bankAccount.Name,
		"bank_code":          bankAccount.Brand,
	}
	sendSuccessResponse(w, response)
}

