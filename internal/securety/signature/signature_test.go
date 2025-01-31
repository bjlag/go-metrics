package signature_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/securety/signature"
)

func TestSignManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := signature.NewSignManager("secret")

		sign := m.Sing([]byte("some data"))
		ok, hash := m.Verify([]byte("some data"), sign)

		assert.True(t, ok)
		assert.Equal(t, sign, hash)
	})

	t.Run("disable", func(t *testing.T) {
		m := signature.NewSignManager("")

		sign := m.Sing([]byte("some data"))
		assert.Empty(t, sign)

		ok, hash := m.Verify([]byte("some data"), sign)
		assert.False(t, ok)
		assert.Empty(t, hash)
	})
}
