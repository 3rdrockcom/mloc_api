package enkrypt

import (
	EVAULT "github.com/epointpayment/mloc_api_go/app/services/enkrypt/client/evault"
	KMS "github.com/epointpayment/mloc_api_go/app/services/enkrypt/client/kms"
)

const (
	// DefaultKeyLength is the default key length required for encryption (AES)
	DefaultKeyLength = 32 // AES-256
)

// Service contains information about enkrypt service
type Service struct {
	kms    *KMS.Client
	evault *EVAULT.Client
}

// New creates an instance of the enkrypt service
func New(cfgKMS KMS.Config, cfgEVAULT EVAULT.Config) (service *Service, err error) {
	// Initialize kms client
	kms, err := KMS.New(cfgKMS)
	if err != nil {
		return
	}

	// Initialize evault client
	evault, err := EVAULT.New(cfgEVAULT)
	if err != nil {
		return
	}

	service = &Service{
		kms:    kms,
		evault: evault,
	}
	return
}

// Get decrypts data using a kms key and returns decrypted data
func (srv Service) Get(keyID, entryID int) (data string, err error) {
	// Get secret key from kms
	key, err := srv.kms.GetKey(keyID)
	if err != nil {
		return
	}

	// Get encrypted data entry from evault
	entry, err := srv.evault.GetEntry(entryID)
	if err != nil {
		return
	}

	// Initialize crypt
	c, err := NewCrypt(key.Key)
	if err != nil {
		return
	}

	// Decrypt data
	data, err = c.Decrypt(entry.Value)
	if err != nil {
		return
	}

	return
}

// Store encrypts data using a kms key and stores it on evault
func (srv Service) Store(data string) (keyID, entryID int, err error) {
	// Generate secret key from kms
	key, err := srv.kms.GenerateKey(DefaultKeyLength)
	if err != nil {
		return
	}

	// Initialize crypt
	c, err := NewCrypt(key.Key)
	if err != nil {
		return
	}

	// Encrypt data
	enc, err := c.Encrypt(data)
	if err != nil {
		return
	}

	// Send encrypted data to evault for storage
	entry, err := srv.evault.StoreEntry(enc)
	if err != nil {
		return
	}

	keyID = key.ID
	entryID = entry.ID
	return
}
