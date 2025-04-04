package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"used2book-backend/internal/services"
	"used2book-backend/internal/models"
	"fmt"
	"time"

	"github.com/omise/omise-go"
	"github.com/streadway/amqp"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	OmiseService *services.OmiseService
	UserService  *services.UserService
	RabbitMQConn *amqp.Connection
}

// ChargeRequest represents the JSON request body structure
type ChargeRequest struct {
	ListingID int `json:"listing_id"`
	BuyerID   int `json:"buyer_id"`
	OfferID   int `json:"offer_id,omitempty"`
}

// ChargeHandler processes a PromptPay payment from the buyer to the seller
// handler/payment_handler.go
// handler/payment_handler.go
func (ph *PaymentHandler) ChargeHandler(w http.ResponseWriter, r *http.Request) {
    var req ChargeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
        return
    }
    log.Println("BuyerID:", req.BuyerID, "ListingID:", req.ListingID, "OfferID:", req.OfferID)

    var amount float64
    var listing *models.ListingDetails
    var sellerID int
    var offerID *int


	listing, err := ph.UserService.GetListingByID(r.Context(), req.ListingID)
        if err != nil || listing == nil {
            sendErrorResponse(w, http.StatusNotFound, "Listing not found")
            return
        }
    if req.OfferID != 0 {
        offer, err := ph.UserService.GetOfferByID(r.Context(), req.OfferID)
        if err != nil {
            sendErrorResponse(w, http.StatusNotFound, "Offer not found")
            return
        }
        if offer.BuyerID != req.BuyerID {
            sendErrorResponse(w, http.StatusForbidden, "You are not the buyer of this offer")
            return
        }
        if offer.Status != "accepted" {
            sendErrorResponse(w, http.StatusBadRequest, "Offer is not accepted")
            return
        }
        amount = offer.OfferedPrice
        // listing, err = ph.UserService.GetListingByID(r.Context(), offer.ListingID)
        // if err != nil || listing == nil {
        //     sendErrorResponse(w, http.StatusNotFound, "Listing not found")
        //     return
        // }
        success, err := ph.UserService.ReserveListingForOffer(r.Context(), offer.ListingID, offer.BuyerID)
        if err != nil || !success {
            log.Println("❌ ReserveListingForOffer Error:", err)
            sendErrorResponse(w, http.StatusInternalServerError, "Failed to reserve listing")
            return
        }
        sellerID = offer.SellerID
        tmpOfferID := req.OfferID
        offerID = &tmpOfferID
    } else {
        // listing, err := ph.UserService.GetListingByID(r.Context(), req.ListingID)
        // if err != nil || listing == nil {
        //     sendErrorResponse(w, http.StatusNotFound, "Listing not found")
        //     return
        // }
        log.Println("listing id :", listing.ListingID)
        log.Println("listing omise:", listing.SellerOmiseID.String)

        amount = float64(listing.Price)
        sellerID = listing.SellerID
        log.Println("listing id :", listing.ListingID)

        success, err := ph.UserService.ReserveListing(r.Context(), req.ListingID, req.BuyerID)
        if err != nil {
            log.Println("❌ ReserveListing Error:", err)
            sendErrorResponse(w, http.StatusInternalServerError, "Failed to process purchase")
            return
        }
        if !success {
            isReserved, isExpired, err := ph.UserService.IsListingReserved(r.Context(), req.ListingID)
            if err != nil {
                log.Println("❌ Error checking reservation status:", err)
                sendErrorResponse(w, http.StatusInternalServerError, "Failed to check listing status")
                return
            }
            if isReserved && !isExpired {
                sendErrorResponse(w, http.StatusConflict, "This book is currently reserved by another buyer.")
                return
            } else if isReserved && isExpired {
                if err := ph.UserService.ExpireReservedListing(r.Context(), req.ListingID); err != nil {
                    log.Println("❌ ExpireReservedListing Error:", err)
                    sendErrorResponse(w, http.StatusInternalServerError, "Failed to process listing status")
                    return
                }
                success, err = ph.UserService.ReserveListing(r.Context(), req.ListingID, req.BuyerID)
                if err != nil || !success {
                    sendErrorResponse(w, http.StatusConflict, "Book is no longer available")
                    return
                }
            } else {
                sendErrorResponse(w, http.StatusConflict, "Book is not available for sale")
                return
            }
        }
        offerID = nil

        if listing == nil {
            log.Println("❌ Listing became nil after reservation")
            sendErrorResponse(w, http.StatusInternalServerError, "Internal server error: listing data missing")
            return
        }
    }


    if !listing.SellerOmiseID.Valid {
        _ = ph.UserService.ExpireReservedListing(r.Context(), req.ListingID)
        sendErrorResponse(w, http.StatusBadRequest, "Seller does not have a valid Omise account")
        return
    }
    sellerOmiseID := listing.SellerOmiseID.String
    log.Println("seller omise id:", sellerOmiseID)

    charge, err := ph.OmiseService.CreatePromptPayCharge(int64(amount*100), req.ListingID, sellerOmiseID, req.BuyerID, 2, offerID)
    if err != nil {
        _ = ph.UserService.ExpireReservedListing(r.Context(), req.ListingID)
        log.Println("❌ PromptPay Charge Error:", err)
        sendErrorResponse(w, http.StatusInternalServerError, "Payment initiation failed")
        return
    }

    if charge.Source == nil || charge.Source.ScannableCode == nil || charge.Source.ScannableCode.Image == nil {
        _ = ph.UserService.ExpireReservedListing(r.Context(), req.ListingID)
        log.Println("❌ Charge response missing QR code data")
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate QR code")
        return
    }

    transactionStatus := "reserved"
    if req.OfferID != 0 {
        transactionStatus = "offer_accepted"
    }
    err = ph.UserService.CreateTransaction(r.Context(), req.BuyerID, sellerID, req.ListingID, amount, transactionStatus)
    if err != nil {
        _ = ph.UserService.ExpireReservedListing(r.Context(), req.ListingID)
        if strings.Contains(err.Error(), "Duplicate entry") {
            sendErrorResponse(w, http.StatusConflict, "Transaction already initiated")
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
        "message":    "Scan this QR code within 2 minutes to pay.",
    }
    sendSuccessResponse(w, response)
}
// WebhookHandler processes Omise payment confirmations (included for completeness)
func (ph *PaymentHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var event struct {
		Key  string        `json:"key"`
		Data *omise.Charge `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Println("❌ Webhook Decode Error:", err)
		http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
		return
	}

	if event.Key == "charge.complete" && event.Data.Paid {
		log.Println("event.Key:", event.Key)
		log.Println("event.Data.Paid:", event.Data.Paid)

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

		err := ph.UserService.UpdateTransactionStatus(context.Background(), int(listingID), "completed")
		if err != nil {
			log.Println("❌ UpdateTransactionStatus Error:", err)
			http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
			return
		}

		err = ph.UserService.MarkListingAsSold(context.Background(), int(listingID), int(buyerID), amount)
		if err != nil {
			log.Println("❌ MarkListingAsSold Error:", err)
			http.Error(w, "Failed to update listing", http.StatusInternalServerError)
			return
		}

		log.Printf("Payment confirmed for listing %d by buyer %d", int(listingID), int(buyerID))

		// Publish to RabbitMQ
        ch, err := ph.RabbitMQConn.Channel()
        if err != nil {
            log.Println("❌ RabbitMQ Channel Error:", err)
            return // Don’t fail webhook response
        }
        defer ch.Close()

        q, err := ch.QueueDeclare(
            "payment_queue", // New queue for payments
            true, false, false, false, nil,
        )
        if err != nil {
            log.Println("❌ Queue Declare Error:", err)
            return
        }

        noti := map[string]interface{}{
            "user_id":    int(buyerID),
            "type":       "payment_success",
            "message":    fmt.Sprintf("Payment succeeded!"),
            "related_id": event.Data.ID, // charge_id
            "created_at": time.Now(),
        }
        body, _ := json.Marshal(noti)
        err = ch.Publish(
            "", q.Name, false, false,
            amqp.Publishing{
                ContentType: "application/json",
                Body:        body,
            },
        )
        if err != nil {
            log.Println("❌ Publish Error:", err)
        }
		

	} else if event.Key == "charge.expired" {
		listingID, ok := event.Data.Metadata["listing_id"].(float64)
		if !ok {
			log.Println("❌ Invalid listing_id in metadata")
			http.Error(w, "Invalid metadata", http.StatusBadRequest)
			return
		}
		offerID, offerExists := event.Data.Metadata["offer_id"].(float64)
		if offerExists {
			err := ph.UserService.RevertOfferReservation(context.Background(), int(listingID), int(offerID))
			if err != nil {
				log.Println("❌ RevertOfferReservation Error:", err)
				http.Error(w, "Failed to revert offer reservation", http.StatusInternalServerError)
				return
			}
			log.Printf("Charge expired for offer %d, listing %d reverted to for_sale", int(offerID), int(listingID))
		} else {
			err := ph.UserService.ExpireReservedListing(context.Background(), int(listingID))
			if err != nil {
				log.Println("❌ ExpireReservedListing Error:", err)
				http.Error(w, "Failed to expire listing", http.StatusInternalServerError)
				return
			}
			log.Printf("Charge expired for listing %d", int(listingID))
		}
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

