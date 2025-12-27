package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncryptor(t *testing.T) {
	t.Run("Valid 32-byte key", func(t *testing.T) {
		key := make([]byte, 32)
		encryptor, err := NewEncryptor(key)
		assert.NoError(t, err)
		assert.NotNil(t, encryptor)
	})

	t.Run("Invalid key length", func(t *testing.T) {
		key := make([]byte, 16) // Only 16 bytes, need 32
		encryptor, err := NewEncryptor(key)
		assert.Error(t, err)
		assert.Nil(t, encryptor)
		assert.Equal(t, "encryption key must be 32 bytes for AES-256", err.Error())
	})
}

func TestEncryptor_Encrypt(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	t.Run("Encrypt non-empty plaintext", func(t *testing.T) {
		plaintext := "This is a secret message"
		ciphertext, err := encryptor.Encrypt(plaintext)
		assert.NoError(t, err)
		assert.NotEmpty(t, ciphertext)
		assert.NotEqual(t, plaintext, ciphertext)
	})

	t.Run("Encrypt empty plaintext", func(t *testing.T) {
		plaintext := ""
		ciphertext, err := encryptor.Encrypt(plaintext)
		assert.Error(t, err)
		assert.Empty(t, ciphertext)
		assert.Equal(t, "plaintext cannot be empty", err.Error())
	})

	t.Run("Different encryptions of same plaintext produce different ciphertexts", func(t *testing.T) {
		plaintext := "Same message"
		ciphertext1, err := encryptor.Encrypt(plaintext)
		require.NoError(t, err)
		ciphertext2, err := encryptor.Encrypt(plaintext)
		require.NoError(t, err)
		// Due to random nonce, ciphertexts should be different
		assert.NotEqual(t, ciphertext1, ciphertext2)
	})
}

func TestEncryptor_Decrypt(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	t.Run("Decrypt valid ciphertext", func(t *testing.T) {
		plaintext := "Secret data to encrypt and decrypt"
		ciphertext, err := encryptor.Encrypt(plaintext)
		require.NoError(t, err)

		decrypted, err := encryptor.Decrypt(ciphertext)
		assert.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("Decrypt empty ciphertext", func(t *testing.T) {
		ciphertext := ""
		decrypted, err := encryptor.Decrypt(ciphertext)
		assert.Error(t, err)
		assert.Empty(t, decrypted)
		assert.Equal(t, "ciphertext cannot be empty", err.Error())
	})

	t.Run("Decrypt invalid ciphertext", func(t *testing.T) {
		ciphertext := "not-a-valid-ciphertext"
		decrypted, err := encryptor.Decrypt(ciphertext)
		assert.Error(t, err)
		assert.Empty(t, decrypted)
	})

	t.Run("Decrypt with wrong key", func(t *testing.T) {
		plaintext := "Original message"
		ciphertext, err := encryptor.Encrypt(plaintext)
		require.NoError(t, err)

		// Create a different encryptor with a different key
		wrongKey := make([]byte, 32)
		for i := range wrongKey {
			wrongKey[i] = byte(255 - i)
		}
		wrongEncryptor, err := NewEncryptor(wrongKey)
		require.NoError(t, err)

		// Try to decrypt with wrong key
		decrypted, err := wrongEncryptor.Decrypt(ciphertext)
		assert.Error(t, err)
		assert.Empty(t, decrypted)
	})
}

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i * 7 % 256)
	}
	encryptor, err := NewEncryptor(key)
	require.NoError(t, err)

	testCases := []string{
		"Short",
		"A longer message with special characters: !@#$%^&*()",
		"Unicode characters: 你好世界 🚀 ñáéíóú",
		"Numbers and symbols: 1234567890 -=[]\\;',./",
		"Multi\nline\ntext\nwith\nnewlines",
	}

	for _, tc := range testCases {
		t.Run("Roundtrip: "+tc[:min(len(tc), 20)], func(t *testing.T) {
			ciphertext, err := encryptor.Encrypt(tc)
			require.NoError(t, err)
			assert.NotEmpty(t, ciphertext)

			decrypted, err := encryptor.Decrypt(ciphertext)
			require.NoError(t, err)
			assert.Equal(t, tc, decrypted)
		})
	}
}

func TestGenerateKey(t *testing.T) {
	t.Run("Generate key of correct length", func(t *testing.T) {
		key, err := GenerateKey()
		assert.NoError(t, err)
		assert.NotNil(t, key)
		assert.Equal(t, 32, len(key))
	})

	t.Run("Generated keys are different", func(t *testing.T) {
		key1, err := GenerateKey()
		require.NoError(t, err)
		key2, err := GenerateKey()
		require.NoError(t, err)
		assert.NotEqual(t, key1, key2)
	})

	t.Run("Generated key works with encryptor", func(t *testing.T) {
		key, err := GenerateKey()
		require.NoError(t, err)

		encryptor, err := NewEncryptor(key)
		require.NoError(t, err)

		plaintext := "Test message with generated key"
		ciphertext, err := encryptor.Encrypt(plaintext)
		require.NoError(t, err)

		decrypted, err := encryptor.Decrypt(ciphertext)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
