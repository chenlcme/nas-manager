package crypto

import (
	"testing"
)

func TestCrypto_GenerateSalt(t *testing.T) {
	c := NewCrypto()

	salt1, err := c.GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	if len(salt1) != saltLen {
		t.Errorf("Expected salt length %d, got %d", saltLen, len(salt1))
	}

	// Generate another salt and verify it's different
	salt2, err := c.GenerateSalt()
	if err != nil {
		t.Fatalf("Failed to generate second salt: %v", err)
	}

	if string(salt1) == string(salt2) {
		t.Error("Expected different salts, got same")
	}
}

func TestCrypto_DeriveKey(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")

	key := c.DeriveKey(password, salt)

	if len(key) != keyLen {
		t.Errorf("Expected key length %d, got %d", keyLen, len(key))
	}

	// Same password and salt should produce same key
	key2 := c.DeriveKey(password, salt)
	if string(key) != string(key2) {
		t.Error("Expected same key for same password and salt")
	}

	// Different password should produce different key
	key3 := c.DeriveKey("differentpassword", salt)
	if string(key) == string(key3) {
		t.Error("Expected different key for different password")
	}
}

func TestCrypto_EncryptDecrypt(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")
	key := c.DeriveKey(password, salt)

	plaintext := []byte("Hello, World!")

	ciphertext, err := c.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := c.Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestCrypto_EncryptDecrypt_Base64(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")
	key := c.DeriveKey(password, salt)

	plaintext := []byte("Hello, World! 你好世界！")

	encoded, err := c.EncryptToBase64(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt to base64: %v", err)
	}

	decrypted, err := c.DecryptFromBase64(encoded, key)
	if err != nil {
		t.Fatalf("Failed to decrypt from base64: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestCrypto_DecryptWithWrongKey(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")
	key := c.DeriveKey(password, salt)

	plaintext := []byte("Hello, World!")

	ciphertext, err := c.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Try to decrypt with wrong key
	wrongKey := c.DeriveKey("wrongpassword", salt)
	_, err = c.Decrypt(ciphertext, wrongKey)
	if err == nil {
		t.Error("Expected error when decrypting with wrong key")
	}
}

func TestCrypto_GenerateVerifyValue(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")

	verifyValue, err := c.GenerateVerifyValue(password, salt)
	if err != nil {
		t.Fatalf("Failed to generate verify value: %v", err)
	}

	if verifyValue == "" {
		t.Error("Expected non-empty verify value")
	}

	// Verify value should be verifiable with same password and salt
	if !c.VerifyPassword(password, salt, verifyValue) {
		t.Error("Expected password to verify against generated verify value")
	}
}

func TestCrypto_VerifyPassword(t *testing.T) {
	c := NewCrypto()
	password := "testpassword123"
	salt := []byte("testsalt123456789012345678901234")

	verifyValue, err := c.GenerateVerifyValue(password, salt)
	if err != nil {
		t.Fatalf("Failed to generate verify value: %v", err)
	}

	// Correct password should verify
	if !c.VerifyPassword(password, salt, verifyValue) {
		t.Error("Expected password to verify")
	}

	// Wrong password should not verify
	if c.VerifyPassword("wrongpassword", salt, verifyValue) {
		t.Error("Expected wrong password to not verify")
	}
}
