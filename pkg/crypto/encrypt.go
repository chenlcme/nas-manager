package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// PBKDF2 参数
	pbkdf2Iterations = 100000
	keyLen           = 32 // 256 bits
	saltLen          = 32 // 256 bits

	// 验证字符串，用于验证密码是否正确
	verifyString = "NAS_MANAGER_VERIFY"
)

// Crypto - 加密解密工具
type Crypto struct{}

// NewCrypto - 创建加密工具实例
func NewCrypto() *Crypto {
	return &Crypto{}
}

// GenerateSalt - 生成随机盐值
func (c *Crypto) GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

// DeriveKey - 使用 PBKDF2 派生密钥
func (c *Crypto) DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, pbkdf2Iterations, keyLen, sha256.New)
}

// Encrypt - 使用 ChaCha20-Poly1305 加密数据
func (c *Crypto) Encrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != keyLen {
		return nil, errors.New("invalid key length")
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密并附加认证标签
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt - 使用 ChaCha20-Poly1305 解密数据
func (c *Crypto) Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != keyLen {
		return nil, errors.New("invalid key length")
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	actualCiphertext := ciphertext[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptToBase64 - 加密并返回 Base64 编码
func (c *Crypto) EncryptToBase64(plaintext, key []byte) (string, error) {
	ciphertext, err := c.Encrypt(plaintext, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptFromBase64 - 解密 Base64 编码的数据
func (c *Crypto) DecryptFromBase64(encoded string, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return c.Decrypt(ciphertext, key)
}

// GenerateVerifyValue - 生成验证值，用于验证密码是否正确
func (c *Crypto) GenerateVerifyValue(password string, salt []byte) (string, error) {
	key := c.DeriveKey(password, salt)
	encrypted, err := c.EncryptToBase64([]byte(verifyString), key)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

// VerifyPassword - 验证密码是否正确
func (c *Crypto) VerifyPassword(password string, salt []byte, verifyValue string) bool {
	key := c.DeriveKey(password, salt)
	decrypted, err := c.DecryptFromBase64(verifyValue, key)
	if err != nil {
		return false
	}
	return string(decrypted) == verifyString
}
