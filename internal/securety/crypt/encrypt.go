package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"os"
)

// EncryptManager отвечает за шифрование данных используя публичный ключ.
type EncryptManager struct {
	publicKey *rsa.PublicKey
}

// NewEncryptManager создает менеджер управляющий шифрование.
func NewEncryptManager(pathToPublicKey string) (*EncryptManager, error) {
	if pathToPublicKey == "" {
		return &EncryptManager{}, nil
	}

	publicKey, err := parsePublicKey(pathToPublicKey)
	if err != nil {
		return nil, err
	}

	return &EncryptManager{
		publicKey: publicKey,
	}, nil
}

// Encrypt шифрует переданные данные.
func (c EncryptManager) Encrypt(data []byte) ([]byte, error) {
	if !c.isEnabled() {
		return data, nil
	}

	dataLen := len(data)
	sizeOneMsg := c.publicKey.Size() - 11

	sizeCipherData := int(math.Ceil(float64(dataLen)/float64(sizeOneMsg))) * c.publicKey.Size()
	result := make([]byte, 0, sizeCipherData)

	for start := 0; start < dataLen; start += sizeOneMsg {
		end := start + sizeOneMsg
		if end > dataLen {
			end = dataLen
		}

		cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, data[start:end])
		if err != nil {
			return nil, fmt.Errorf("rsa.EncryptPKCS1v15: %w", err)
		}

		result = append(result, cipherData...)
	}

	return result, nil
}

func (c EncryptManager) isEnabled() bool {
	return c.publicKey != nil
}

func parsePublicKey(pathToKey string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(pathToKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return publicKey, nil
}
