package enkrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt produces a base64-wrapped/AES-encrypted string using AES-GCM
func (c *Crypt) EncryptAESGCM(text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// The nonce needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, gcm.NonceSize()+len(plaintext)+gcm.Overhead())
	nonce := ciphertext[:gcm.NonceSize()]
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	buffer.Write(nonce)
	buffer.Write(gcm.Seal(nil, nonce, plaintext, nil))
	ciphertext = buffer.Bytes()

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-wrapped/AES-encrypted string using AES-GCM
func (c *Crypt) DecryptAESGCM(cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// The nonce needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if byteLen := len(ciphertext); byteLen < gcm.NonceSize() {
		return "", fmt.Errorf("invalid cipher size %d", byteLen)
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	ciphertext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(ciphertext), nil
}
