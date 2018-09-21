package customer

import (
	"github.com/epointpayment/mloc_api_go/app/config"
	"github.com/epointpayment/mloc_api_go/app/services/enkrypt"
	"github.com/epointpayment/mloc_api_go/app/services/enkrypt/client/evault"
	"github.com/epointpayment/mloc_api_go/app/services/enkrypt/client/kms"
)

// getEnkryptConfig get KMS and Evault configuration information for Enkrypt service
func (a *BankAccount) getEnkryptConfig() (kmsConfig kms.Config, evaultConfig evault.Config, err error) {
	cfg := config.Get()

	kmsConfig = kms.Config{
		BaseURL:   cfg.KMS.BaseURL,
		ProgramID: cfg.KMS.ProgramID,
		Username:  cfg.KMS.Username,
		Password:  cfg.KMS.Password,
	}

	evaultConfig = evault.Config{
		BaseURL:     cfg.Evault.BaseURL,
		ProgramID:   cfg.Evault.ProgramID,
		PartitionID: cfg.Evault.PartitionID,
		Username:    cfg.Evault.Username,
		Password:    cfg.Evault.Password,
	}

	return
}

// encrypt encrypts text using KMS key and stores it on Evault using Enkrypt service
func (a *BankAccount) encrypt(text string) (kmsID int, evaultID int, err error) {
	// Get enkrypt configuration
	kmsConfig, evaultConfig, err := a.getEnkryptConfig()
	if err != nil {
		return
	}

	// Initialize enkrypt service
	krypt, err := enkrypt.New(kmsConfig, evaultConfig)
	if err != nil {
		return
	}

	// Store encrypted text
	kmsID, evaultID, err = krypt.Store(text)
	if err != nil {
		return
	}

	return
}

// decrypt gets encrypted text from Evault and decrypts it using a KMS key using Enkrypt service
func (a *BankAccount) decrypt(kmsID, evaultID int) (text string, err error) {
	// Get enkrypt configuration
	kmsConfig, evaultConfig, err := a.getEnkryptConfig()
	if err != nil {
		return
	}

	// Initialize enkrypt service
	krypt, err := enkrypt.New(kmsConfig, evaultConfig)
	if err != nil {
		return
	}

	// Get decrypted text
	text, err = krypt.Get(kmsID, evaultID)
	if err != nil {
		return
	}

	return
}
