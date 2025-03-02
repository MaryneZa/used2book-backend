package services

import (
	"fmt"
	"log"
	"os"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// OmiseService handles Omise API interactions
type OmiseService struct {
	Client *omise.Client
}

// NewOmiseService initializes Omise with your API keys
func NewOmiseService() *OmiseService {
	client, err := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	if err != nil {
		log.Fatal("Omise client error:", err)
	}
	return &OmiseService{Client: client}
}

// CreateCharge processes a payment from a buyer
func (o *OmiseService) CreateCharge(amount int64, token string, sellerRecipientID string) (*omise.Charge, error) {
	charge := &omise.Charge{}

	// Create a charge using the buyer's card token
	createCharge := &operations.CreateCharge{
		Amount:      amount, // Amount in satangs (e.g., 1000 satangs = 10 THB)
		Currency:    "THB",
		Card:        token, // Token from frontend (e.g., "tok_xxxx")
		Description: "Book purchase",
		Metadata: map[string]interface{}{
			"seller_id": sellerRecipientID, // Store seller's Omise Recipient ID for tracking
		},
	}

	// Corrected call to Client.Do (no context)
	err := o.Client.Do(charge, createCharge)
	if err != nil {
		return nil, fmt.Errorf("failed to create charge: %v", err)
	}

	return charge, nil
}

func (o *OmiseService) CreateRecipient(bankAccountNumber, bankAccountName, bankCode string) (string, error) {
	recipient := &omise.Recipient{}

	// Create a recipient for the seller
	createRecipient := &operations.CreateRecipient{
		Name: bankAccountName,
		Type: "individual", // or "corporation" for businesses
		BankAccount: &omise.BankAccount{
			Number: bankAccountNumber,
			Name:   bankAccountName,
			Brand:  bankCode, // Example: "bbl" for Bangkok Bank, "scb" for Siam Commercial Bank
		},
	}

	// Execute the API call
	err := o.Client.Do(recipient, createRecipient)
	if err != nil {
		return "", fmt.Errorf("failed to create recipient: %v", err)
	}

	return recipient.ID, nil
}