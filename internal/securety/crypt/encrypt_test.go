package crypt_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/securety/crypt"
)

func TestEncryptManager_Encrypt(t *testing.T) {
	publicKeyPEM := `-----BEGIN RSA PUBLIC KEY-----
MBgCEQDXosrBcPG1oos0rB7pkKpFAgMBAAE=
-----END RSA PUBLIC KEY-----`

	privateKeyPEM := `-----BEGIN RSA PRIVATE KEY-----
MGICAQACEQDXosrBcPG1oos0rB7pkKpFAgMBAAECECvZh9+kZxKnNpQiEvvF7LUC
CQDul6leSB/I3wIJAOdeW3U83R1bAggk17aHoIuH8QIJAJiC2katmBOBAgg4SFxm
/kd1kQ==
-----END RSA PRIVATE KEY-----`

	tmpFilePublicKey, err := os.CreateTemp("", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFilePublicKey.Name())
	}()

	if _, err = tmpFilePublicKey.Write([]byte(publicKeyPEM)); err != nil {
		t.Fatal(err)
	}

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

	t.Run("encrypt", func(t *testing.T) {
		encrypt, err := crypt.NewEncryptManager(tmpFilePublicKey.Name())
		if err != nil {
			t.Fatal(err)
		}

		decrypt, err := crypt.NewDecryptManager(tmpFilePrivateKey.Name())
		if err != nil {
			t.Fatal(err)
		}

		cipherData, err := encrypt.Encrypt([]byte("Some metrics"))
		assert.NoError(t, err)

		decryptedData, err := decrypt.Decrypt(cipherData)
		assert.NoError(t, err)

		assert.Equal(t, []byte("Some metrics"), decryptedData)
	})

	t.Run("disable", func(t *testing.T) {
		encrypt, err := crypt.NewEncryptManager("")
		assert.NoError(t, err)

		cipherData, err := encrypt.Encrypt([]byte("Some metrics"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("Some metrics"), cipherData)
	})
}
