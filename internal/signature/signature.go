package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SignManager struct {
	secretKey []byte
	enable    bool
}

func NewSignManager(secretKey string) *SignManager {
	return &SignManager{
		secretKey: []byte(secretKey),
		enable:    len(secretKey) > 0,
	}
}

func (m SignManager) Sing(data []byte) string {
	if !m.enable {
		return ""
	}

	return hex.EncodeToString(m.new(data))
}

func (m SignManager) Verify(data []byte, signature string) (bool, string) {
	if !m.enable {
		return false, ""
	}

	sign, err := hex.DecodeString(signature)
	if err != nil {
		return false, ""
	}

	dataSign := m.new(data)

	return hmac.Equal(dataSign, sign), hex.EncodeToString(dataSign)
}

func (m SignManager) Enable() bool {
	return m.enable
}

func (m SignManager) new(data []byte) []byte {
	h := hmac.New(sha256.New, m.secretKey)
	h.Write(data)
	return h.Sum(nil)
}
