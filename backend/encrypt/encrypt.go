package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/url"

	"golang.org/x/crypto/pbkdf2"
)

const (
	password   = "password"
	iterations = 2000
	keyLen     = 16 // AES-128
	ivLen      = 16
	saltLen    = 8
)

// Encrypt encrypts plain text using AES-128-CBC with PBKDF2 key derivation.
// Returns base64-encoded ciphertext and base64-encoded salt.
// Compatible with the Ruby OpenSSL implementation.
func Encrypt(plainText string) (string, string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	keyIV := pbkdf2.Key([]byte(password), salt, iterations, keyLen+ivLen, sha256.New)
	key := keyIV[:keyLen]
	iv := keyIV[keyLen:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	// PKCS7 padding
	plainBytes := []byte(plainText)
	padding := aes.BlockSize - len(plainBytes)%aes.BlockSize
	for i := 0; i < padding; i++ {
		plainBytes = append(plainBytes, byte(padding))
	}

	encrypted := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, plainBytes)

	encText := base64.StdEncoding.EncodeToString(encrypted)
	encSalt := base64.StdEncoding.EncodeToString(salt)

	return encText, encSalt, nil
}

// Decrypt decrypts base64-encoded ciphertext using the given base64-encoded salt.
// Compatible with the Ruby OpenSSL implementation.
func Decrypt(encryptedText, saltStr string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return "", err
	}

	keyIV := pbkdf2.Key([]byte(password), salt, iterations, keyLen+ivLen, sha256.New)
	key := keyIV[:keyLen]
	iv := keyIV[keyLen:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove PKCS7 padding
	if len(decrypted) > 0 {
		padding := int(decrypted[len(decrypted)-1])
		if padding > 0 && padding <= aes.BlockSize {
			decrypted = decrypted[:len(decrypted)-padding]
		}
	}

	// CGI.unescape equivalent
	result, err := url.QueryUnescape(string(decrypted))
	if err != nil {
		// If unescape fails, return as-is (data may not be URL-encoded)
		return string(decrypted), nil
	}
	return result, nil
}
