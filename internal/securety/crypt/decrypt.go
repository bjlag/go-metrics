package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// DecryptManager отвечает за расшифровку данных используя приватный ключ.
type DecryptManager struct {
	privateKey *rsa.PrivateKey
}

// NewDecryptManager создает менеджер управляющий расшифровкой.
func NewDecryptManager(pathToPrivateKey string) (*DecryptManager, error) {
	if pathToPrivateKey == "" {
		return &DecryptManager{}, nil
	}

	privateKey, err := parsePrivateKey(pathToPrivateKey)
	if err != nil {
		return nil, err
	}

	return &DecryptManager{
		privateKey: privateKey,
	}, nil
}

// Decrypt расшифровывает переданные данные.
func (c DecryptManager) Decrypt(cipherData []byte) ([]byte, error) {
	if !c.isEnabled() {
		return cipherData, nil
	}

	dataLen := len(cipherData)
	sizeOneMsg := c.privateKey.PublicKey.Size()

	result := make([]byte, 0, dataLen)

	for start := 0; start < dataLen; start += sizeOneMsg {
		end := start + sizeOneMsg
		if end > dataLen {
			end = dataLen
		}

		data, err := rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, cipherData[start:end])
		if err != nil {
			return nil, fmt.Errorf("rsa.DecryptOAEP: %w", err)
		}

		result = append(result, data...)
	}

	return result, nil
}

func (c DecryptManager) isEnabled() bool {
	return c.privateKey != nil
}

func parsePrivateKey(pathToPrivateKey string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(pathToPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
