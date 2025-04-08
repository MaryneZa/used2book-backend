package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"used2book-backend/internal/models"
	"used2book-backend/internal/services"
)

type PaymentHandler struct {
	UserService  *services.UserService
	RabbitMQConn *amqp.Connection
}

// CheckoutRequest represents the JSON request body structure
type CheckoutRequest struct {
	ListingID int `json:"listing_id"`
	BuyerID   int `json:"buyer_id"`
	OfferID   int `json:"offer_id,omitempty"`
}

func (ph *PaymentHandler) CheckOutHandler(w http.ResponseWriter, r *http.Request) {
	// reserve listing at = stripe expired_at

	if err := godotenv.Load(); err != nil {
		log.Println(errors.New("failed to load stripe_sk_key .env file"))

	}
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}
	log.Println("BuyerID:", req.BuyerID, "ListingID:", req.ListingID, "OfferID:", req.OfferID)

	var amount float64
	var listing *models.ListingDetails
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

		success, err := ph.UserService.ReserveListing(r.Context(), offer.ListingID, offer.BuyerID)
		if err != nil || !success {
			log.Println("❌ ReserveListingForOffer Error:", err)
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to reserve listing")
			return
		}
		tmpOfferID := req.OfferID
		offerID = &tmpOfferID
	} else {

		log.Println("no-offer listing id :", listing.ListingID)

		amount = float64(listing.Price)
		offerID = nil

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
	}

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:3000/user/account/purchase"),
		CancelURL:  stripe.String("http://localhost:3000/user/cancel"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("thb"),
					UnitAmount: stripe.Int64(int64(amount * 100)), // THB in satangs
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(listing.Title),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		ExpiresAt: stripe.Int64(time.Now().Unix() + 1800), // 5 minutes
		Metadata: map[string]string{
			"listing_id": fmt.Sprintf("%d", req.ListingID),
			"buyer_id":   fmt.Sprintf("%d", req.BuyerID),
		},
	}
	if offerID != nil {
		params.Metadata["offer_id"] = fmt.Sprintf("%d", *offerID)
	}

	checkout_session, err := session.New(params)

	log.Println("checkout_session:", checkout_session)

	if err != nil {
		log.Println("Stripe session error:", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Unable to create payment session")
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success":          true,
		"checkout_url":     checkout_session.URL,
		"session_id":       checkout_session.ID,
		"checkout_session": checkout_session,
	})
}

func (ph *PaymentHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		return
	}
	log.Println("📦 Raw payload:", string(payload))

	sigHeader := r.Header.Get("Stripe-Signature")

	if err := godotenv.Load(); err != nil {
		log.Println(errors.New("failed to load stripe_sk_key .env file"))

	}
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET") // from your Stripe dashboard

	// event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)

	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, endpointSecret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})

	if err != nil {
		log.Printf("⚠️  Webhook signature verification failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Optional: log event for debugging
	log.Printf("✅✅ Received event: %s", event.Type)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("⚠️  Failed to parse session: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// ✅ Access metadata (listing_id, buyer_id, offer_id) from your session
		amount_total := float64(session.AmountTotal / 100)
		listingIDStr := session.Metadata["listing_id"]
		buyerIDStr := session.Metadata["buyer_id"]
		var offerID *int
		if offerIDStr, ok := session.Metadata["offer_id"]; ok && offerIDStr != "" {
			parsedID, err := strconv.Atoi(offerIDStr)
			if err == nil {
				offerID = &parsedID
			}
		}

		listingID, err := strconv.Atoi(listingIDStr) // Convert to int
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid listing ID")
			return
		}

		buyerID, err := strconv.Atoi(buyerIDStr) // Convert to int
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid buyer ID")
			return
		}

		sessionID := session.ID
		log.Println("💳 Stripe Session ID:", sessionID)

		err = ph.UserService.CreateTransaction(context.Background(), sessionID, buyerID, listingID, offerID, amount_total, "completed")
		if err != nil {
			log.Println("❌ UpdateTransactionStatus Error:", err)
			http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
			return
		}

		err = ph.UserService.MarkListingAsSold(context.Background(), listingID, buyerID)
		if err != nil {
			log.Println("❌ MarkListingAsSold Error:", err)
			http.Error(w, "Failed to update listing", http.StatusInternalServerError)
			return
		}

		log.Printf("Payment confirmed for listing %d by buyer %d", listingID, buyerID)

		err = ph.UserService.RemoveFromCart(r.Context(), buyerID, listingID)
		if err != nil {
			log.Println("❌ Wishlist Error:", err)
			// sendErrorResponse(w, http.StatusConflict, "Wishlist error: "+err.Error())
			// return
		}

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

        listing, err := ph.UserService.GetListingByID(r.Context(), listingID)
	if err != nil || listing == nil {
		sendErrorResponse(w, http.StatusNotFound, "Listing not found")
		return
	}

		noti := map[string]interface{}{
			"buyer_id":    int(buyerID),
            "listing_id":   listingID,
            "seller_id":    listing.SellerID,
			"type":       "payment_success",
			"message":    "Payment succeeded!",
			"related_id": sessionID, // charge_id
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

		log.Printf("💰 Payment success! Listing ID: %d, Buyer ID: %d, Offer ID: %d", listingID, buyerID, offerID)

	case "checkout.session.expired":
		// Optional: handle expiration
		log.Println("🕒 Checkout session expired")

	default:
		log.Println("Unhandled event type:", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// // WebhookHandler processes Omise payment confirmations (included for completeness)
// func (ph *PaymentHandler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
// 	var event struct {
// 		Key  string        `json:"key"`
// 		Data *omise.Charge `json:"data"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
// 		log.Println("❌ Webhook Decode Error:", err)
// 		http.Error(w, "Invalid webhook payload", http.StatusBadRequest)
// 		return
// 	}

// 	if event.Key == "charge.complete" && event.Data.Paid {
// 		log.Println("event.Key:", event.Key)
// 		log.Println("event.Data.Paid:", event.Data.Paid)

// 		listingID, ok := event.Data.Metadata["listing_id"].(float64)
// 		if !ok {
// 			log.Println("❌ Invalid listing_id in metadata")
// 			http.Error(w, "Invalid metadata", http.StatusBadRequest)
// 			return
// 		}
// 		buyerID, ok := event.Data.Metadata["buyer_id"].(float64)
// 		if !ok {
// 			log.Println("❌ BuyerID not found in metadata")
// 			http.Error(w, "Invalid metadata", http.StatusBadRequest)
// 			return
// 		}
// 		amount := float64(event.Data.Amount / 100)

// 		err := ph.UserService.UpdateTransactionStatus(context.Background(), int(listingID), "completed")
// 		if err != nil {
// 			log.Println("❌ UpdateTransactionStatus Error:", err)
// 			http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
// 			return
// 		}

// 		err = ph.UserService.MarkListingAsSold(context.Background(), int(listingID), int(buyerID), amount)
// 		if err != nil {
// 			log.Println("❌ MarkListingAsSold Error:", err)
// 			http.Error(w, "Failed to update listing", http.StatusInternalServerError)
// 			return
// 		}

// 		log.Printf("Payment confirmed for listing %d by buyer %d", int(listingID), int(buyerID))

// 		// Publish to RabbitMQ
//         ch, err := ph.RabbitMQConn.Channel()
//         if err != nil {
//             log.Println("❌ RabbitMQ Channel Error:", err)
//             return // Don’t fail webhook response
//         }
//         defer ch.Close()

//         q, err := ch.QueueDeclare(
//             "payment_queue", // New queue for payments
//             true, false, false, false, nil,
//         )
//         if err != nil {
//             log.Println("❌ Queue Declare Error:", err)
//             return
//         }

//         noti := map[string]interface{}{
//             "user_id":    int(buyerID),
//             "type":       "payment_success",
//             "message":    fmt.Sprintf("Payment succeeded!"),
//             "related_id": event.Data.ID, // charge_id
//             "created_at": time.Now(),
//         }
//         body, _ := json.Marshal(noti)
//         err = ch.Publish(
//             "", q.Name, false, false,
//             amqp.Publishing{
//                 ContentType: "application/json",
//                 Body:        body,
//             },
//         )
//         if err != nil {
//             log.Println("❌ Publish Error:", err)
//         }

// 	} else if event.Key == "charge.expired" {
// 		listingID, ok := event.Data.Metadata["listing_id"].(float64)
// 		if !ok {
// 			log.Println("❌ Invalid listing_id in metadata")
// 			http.Error(w, "Invalid metadata", http.StatusBadRequest)
// 			return
// 		}
// 		offerID, offerExists := event.Data.Metadata["offer_id"].(float64)
// 		if offerExists {
// 			err := ph.UserService.RevertOfferReservation(context.Background(), int(listingID), int(offerID))
// 			if err != nil {
// 				log.Println("❌ RevertOfferReservation Error:", err)
// 				http.Error(w, "Failed to revert offer reservation", http.StatusInternalServerError)
// 				return
// 			}
// 			log.Printf("Charge expired for offer %d, listing %d reverted to for_sale", int(offerID), int(listingID))
// 		} else {
// 			err := ph.UserService.ExpireReservedListing(context.Background(), int(listingID))
// 			if err != nil {
// 				log.Println("❌ ExpireReservedListing Error:", err)
// 				http.Error(w, "Failed to expire listing", http.StatusInternalServerError)
// 				return
// 			}
// 			log.Printf("Charge expired for listing %d", int(listingID))
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
// // CreateOrUpdateOmiseAccountHandler sets up or updates a seller's Omise recipient
// func (ph *PaymentHandler) CreateOrUpdateOmiseAccountHandler(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := r.Context().Value("user_id").(int)
// 	if !ok || userID == 0 {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var req struct {
// 		BankAccountNumber string `json:"bank_account_number"`
// 		BankAccountName   string `json:"bank_account_name"`
// 		BankCode          string `json:"bank_code"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	if req.BankAccountNumber == "" || req.BankAccountName == "" || req.BankCode == "" {
// 		http.Error(w, "Bank details are required", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := ph.UserService.GetMe(r.Context(), userID)
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	var recipientID string
// 	if user.OmiseAccountID.Valid {
// 		_, err := ph.OmiseService.GetRecipient(user.OmiseAccountID.String)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "not_found") {
// 				log.Printf("Recipient %s not found, creating a new one", user.OmiseAccountID.String)
// 				recipientID, err = ph.OmiseService.CreateRecipient(req.BankAccountNumber, req.BankAccountName, req.BankCode)
// 				if err != nil {
// 					log.Println("Error creating Omise recipient:", err)
// 					http.Error(w, "Failed to create Omise account", http.StatusInternalServerError)
// 					return
// 				}
// 			} else {
// 				log.Println("Error retrieving Omise recipient:", err)
// 				http.Error(w, "Failed to verify Omise account", http.StatusInternalServerError)
// 				return
// 			}
// 		} else {
// 			recipientID, err = ph.OmiseService.UpdateRecipient(user.OmiseAccountID.String, req.BankAccountNumber, req.BankAccountName, req.BankCode)
// 			if err != nil {
// 				log.Println("Error updating Omise recipient:", err)
// 				http.Error(w, "Failed to update Omise account", http.StatusInternalServerError)
// 				return
// 			}
// 		}
// 	} else {
// 		recipientID, err = ph.OmiseService.CreateRecipient(req.BankAccountNumber, req.BankAccountName, req.BankCode)
// 		if err != nil {
// 			log.Println("Error creating Omise recipient:", err)
// 			http.Error(w, "Failed to create Omise account", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	err = ph.UserService.UpdateOmiseAccountID(r.Context(), userID, recipientID)
// 	if err != nil {
// 		log.Println("Error saving Omise account ID to DB:", err)
// 		http.Error(w, "Failed to store Omise account ID", http.StatusInternalServerError)
// 		return
// 	}

// 	resp := map[string]interface{}{
// 		"success":         true,
// 		"omise_recipient": recipientID,
// 	}
// 	json.NewEncoder(w).Encode(resp)
// }

// // GetBankAccountInfoHandler retrieves the user's bank account details from Omise
// func (ph *PaymentHandler) GetBankAccountInfoHandler(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := r.Context().Value("user_id").(int)
// 	if !ok || userID == 0 {
// 		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
// 		return
// 	}

// 	user, err := ph.UserService.GetMe(r.Context(), userID)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusNotFound, "User not found")
// 		return
// 	}

// 	if !user.OmiseAccountID.Valid {
// 		sendErrorResponse(w, http.StatusBadRequest, "No bank account linked to this user")
// 		return
// 	}

// 	recipient, err := ph.OmiseService.GetRecipient(user.OmiseAccountID.String)
// 	if err != nil {
// 		log.Println("❌ Error retrieving Omise recipient:", err)
// 		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve bank account info")
// 		return
// 	}

// 	bankAccount := recipient.BankAccount
// 	response := map[string]interface{}{
// 		"success":           true,
// 		"bank_account_number": bankAccount.Number,
// 		"bank_account_name":   bankAccount.Name,
// 		"bank_code":          bankAccount.Brand,
// 	}
// 	sendSuccessResponse(w, response)
// }
