package enkrypt

const (
	// EncryptionAlgorithmAESGCM is the name of the AES-GCM algorithm
	EncryptionAlgorithmAESGCM = "AES-GCM"

	// EncryptionAlgorithmAESCFB is the name of the AES-CFB algorithm
	EncryptionAlgorithmAESCFB = "AES-CFB"

	// DefaultEncryptionAlgorithm is the default algorithm used for encryption
	DefaultEncryptionAlgorithm = EncryptionAlgorithmAESGCM
)

// Crypt contains information used for encryption and decryption operations
type Crypt struct {
	key       []byte
	algorithm string
}

// NewCrypt creates an instance of the crypt service
func NewCrypt(key string) (c *Crypt, err error) {
	c = new(Crypt)

	// Set secret key
	err = c.setKey(key)
	if err != nil {
		return
	}

	// Set encryption standard
	err = c.setEncryptionAlgorithm(DefaultEncryptionAlgorithm)
	if err != nil {
		return
	}

	return
}

// setKey sets the crypto key
func (c *Crypt) setKey(key string) (err error) {
	keyLen := len(key)

	switch keyLen {
	case 16, 24, 32:
		break
	default:
		err = ErrInvalidKeySize
		return
	}

	c.key = []byte(key)
	return
}

// setEncryptionAlgorithmKey sets the crypto algorithm
func (c *Crypt) setEncryptionAlgorithm(algorithm string) (err error) {
	switch algorithm {
	case
		EncryptionAlgorithmAESGCM,
		EncryptionAlgorithmAESCFB:
		break
	default:
		err = ErrInvalidAlgorithm
		return
	}

	c.algorithm = algorithm
	return
}

// Encrypt produces an encrypted string
func (c *Crypt) Encrypt(text string) (cipherText string, err error) {
	switch c.algorithm {
	case EncryptionAlgorithmAESGCM:
		return c.EncryptAESGCM(text)
	case EncryptionAlgorithmAESCFB:
		return c.EncryptAESCFB(text)
	}

	err = ErrInvalidAlgorithm
	return
}

// Decrypt decrypts an encrypted string
func (c *Crypt) Decrypt(cryptoText string) (text string, err error) {
	switch c.algorithm {
	case EncryptionAlgorithmAESGCM:
		return c.DecryptAESGCM(cryptoText)
	case EncryptionAlgorithmAESCFB:
		return c.DecryptAESCFB(cryptoText)
	}

	err = ErrInvalidAlgorithm
	return
}
