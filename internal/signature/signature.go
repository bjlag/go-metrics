package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SignManager struct {
	secretKey []byte
}

func NewSignManager(secretKey string) *SignManager {
	return &SignManager{
		secretKey: []byte(secretKey),
	}
}

func (m SignManager) Sing(data []byte) string {
	return hex.EncodeToString(m.new(data))
}

func (m SignManager) Verify(data []byte, signature string) (bool, string) {
	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false, ""
	}

	dataSign := m.new(data)

	return hmac.Equal(dataSign, sign), hex.EncodeToString(dataSign)
}

func (m SignManager) new(data []byte) []byte {
	h := hmac.New(sha256.New, m.secretKey)
	h.Write(data)
	return h.Sum(nil)
}
