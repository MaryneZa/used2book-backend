package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// OmiseService handles Omise API interactions
type OmiseService struct {
	Client *omise.Client
}

// NewOmiseService initializes Omise with API keys
func NewOmiseService() *OmiseService {
	client, err := omise.NewClient(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	if err != nil {
		log.Fatal("Omise client error:", err)
	}
	return &OmiseService{Client: client}
}

// services/omise_service.go
func (o *OmiseService) CreatePromptPayCharge(amount int64, listingID int, sellerRecipientID string, buyerID int, expiresInMinutes int, offerID *int) (*omise.Charge, error) {
    source := &omise.Source{}
    createSource := &operations.CreateSource{
        Type:     "promptpay",
        Amount:   amount,
        Currency: "THB",
    }
    if err := o.Client.Do(source, createSource); err != nil {
        return nil, fmt.Errorf("failed to create PromptPay source: %v", err)
    }

    charge := &omise.Charge{}
    expiresAt := time.Now().Add(time.Duration(expiresInMinutes) * time.Minute)
    createCharge := &operations.CreateCharge{
        Amount:      amount,
        Currency:    "THB",
        Source:      source.ID,
        Description: fmt.Sprintf("Book purchase for listing %d", listingID),
        ExpiresAt:   &expiresAt,
        Metadata: map[string]interface{}{
            "listing_id": listingID,
            "buyer_id":   buyerID,
        },
    }
    if offerID != nil {
        createCharge.Metadata["offer_id"] = *offerID
    }
    if err := o.Client.Do(charge, createCharge); err != nil {
        return nil, fmt.Errorf("failed to create PromptPay charge: %v", err)
    }
    return charge, nil
}

// GetRecipient fetches a recipient's details from Omise
func (o *OmiseService) GetRecipient(recipientID string) (*omise.Recipient, error) {
	recipient := &omise.Recipient{}
	err := o.Client.Do(recipient, &operations.RetrieveRecipient{RecipientID: recipientID})
	if err != nil {
		return nil, err
	}
	return recipient, nil
}

// CreateRecipient creates a new Omise recipient for the seller
func (o *OmiseService) CreateRecipient(bankAccountNumber, bankAccountName, bankCode string) (string, error) {
	recipient := &omise.Recipient{}
	createRecipient := &operations.CreateRecipient{
		Name: bankAccountName,
		Type: "individual",
		BankAccount: &omise.BankAccount{
			Number: bankAccountNumber,
			Name:   bankAccountName,
			Brand:  bankCode,
		},
	}
	err := o.Client.Do(recipient, createRecipient)
	if err != nil {
		return "", fmt.Errorf("failed to create recipient: %v", err)
	}
	return recipient.ID, nil
}

// UpdateRecipient updates an existing Omise recipient
func (o *OmiseService) UpdateRecipient(recipientID, bankAccountNumber, bankAccountName, bankCode string) (string, error) {
	recipient := &omise.Recipient{}
	updateOp := &operations.UpdateRecipient{
		RecipientID: recipientID,
		Name:        bankAccountName,
		Type:        "individual",
		BankAccount: &omise.BankAccount{
			Number: bankAccountNumber,
			Name:   bankAccountName,
			Brand:  bankCode,
		},
	}
	err := o.Client.Do(recipient, updateOp)
	if err != nil {
		return "", fmt.Errorf("failed to update recipient: %v", err)
	}
	return recipient.ID, nil
}