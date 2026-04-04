package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"server/config"

	"github.com/sony/sonyflake"
)

func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getEncryptionKey derives a 32-byte key from JWT_SIGNING_KEY using SHA256
func getEncryptionKey() []byte {
	// Use JWT_SIGNING_KEY as the base key, hash it to get 32 bytes
	hash := sha256.Sum256([]byte(config.Settings.JWT_SIGNING_KEY))
	return hash[:]
}

// EncryptAPIKey encrypts an API key using AES-256-GCM
func EncryptAPIKey(apiKey string) (string, error) {
	key := getEncryptionKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAPIKey decrypts an API key using AES-256-GCM
func DecryptAPIKey(encryptedKey string) (string, error) {
	key := getEncryptionKey()

	data, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

var (
	sf     *sonyflake.Sonyflake
	sfOnce sync.Once
)

// GenerateSnowID 生成一个雪花 ID
func GenerateSnowID() (uint64, error) {
	init := func() {
		startTime := time.Date(2006, 4, 11, 0, 0, 0, 0, time.UTC) // 自定义起始时间（epoch）

		machineID := func() (uint16, error) {
			return 1, nil // 生产环境中应确保每台机器 ID 唯一
		}

		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			StartTime: startTime,
			MachineID: machineID,
		})

		if sf == nil {
			panic("failed to create Sonyflake instance")
		}
	}
	sfOnce.Do(init)
	return sf.NextID()
}
