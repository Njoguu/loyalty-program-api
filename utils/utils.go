package utils

import (
	"time"
	"crypto/rand"
	random"math/rand"
	"encoding/base64"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Virtual Card Generator
func GenerateCardNumber() string {
	source := random.NewSource(time.Now().UnixNano())
	random := random.New(source)

	var digits = []rune("0123456789")
	cardNumber := make([]rune, 13)

	for i := 0; i < 13; i++ {
		cardNumber[i] = digits[random.Intn(len(digits))]
	}

	return string(cardNumber)
}