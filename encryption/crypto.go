package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	KeyLength   = 32
	SaltLength  = 16
	Iteration   = 5
	Memory      = 512 * 1024
	Parallelism = 8
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil

}

func DeriveKey(masterPassword string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(masterPassword),
		salt,
		Iteration,
		Memory,
		Parallelism,
		KeyLength,
	)
}

func Encrypt(key []byte, plainText []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	cipherText := gcm.Seal(nil, nonce, plainText, nil)

	return nonce, cipherText, nil
}

func Decrypt(key []byte, nonce []byte, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, cipherText, nil)
}

func GenerateRandomPassword(bytes int) (string, error) {
	b := make([]byte, bytes)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
