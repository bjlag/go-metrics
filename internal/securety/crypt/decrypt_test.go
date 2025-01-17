package crypt_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/securety/crypt"
)

func TestDecryptManager_Decrypt(t *testing.T) {
	cipherData := "ZPN9nkk7vUl98+mzkprVzRLTyWXDGK6UNQjtLV+PMoio3OZg5aDlO6+UOaRY0Qqv"

	privateKeyPEM := `-----BEGIN RSA PRIVATE KEY-----
MGICAQACEQDXosrBcPG1oos0rB7pkKpFAgMBAAECECvZh9+kZxKnNpQiEvvF7LUC
CQDul6leSB/I3wIJAOdeW3U83R1bAggk17aHoIuH8QIJAJiC2katmBOBAgg4SFxm
/kd1kQ==
-----END RSA PRIVATE KEY-----`

	tmpFilePrivateKey, err := os.CreateTemp("", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFilePrivateKey.Name())
	}()

	if _, err = tmpFilePrivateKey.Write([]byte(privateKeyPEM)); err != nil {
		t.Fatal(err)
	}

	t.Run("decrypt", func(t *testing.T) {
		decrypt, err := crypt.NewDecryptManager(tmpFilePrivateKey.Name())
		assert.NoError(t, err)

		decoded, err := base64.StdEncoding.DecodeString(cipherData)
		assert.NoError(t, err)

		data, err := decrypt.Decrypt(decoded)
		assert.NoError(t, err)
		assert.Equal(t, []byte("Some metrics"), data)
	})

	t.Run("disable", func(t *testing.T) {
		decrypt, err := crypt.NewDecryptManager("")
		assert.NoError(t, err)

		data, err := decrypt.Decrypt([]byte("Some metrics"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("Some metrics"), data)
	})
}
