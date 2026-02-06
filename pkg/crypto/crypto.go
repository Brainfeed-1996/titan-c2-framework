package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Key should be 32 bytes for AES-256
var Key = []byte("thisisausecretkeyforAES256encryption")[:32]

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}

	b := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher(Key)
	if err != nil {
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	if len(decoded) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decoded, decoded)
	return string(decoded), nil
}
