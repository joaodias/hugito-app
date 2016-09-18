package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// AreStringsEqual returns a boolean representing if strings are equal or not.
// If strings are equal true is returned. For all other cases false is
// returned.
func AreStringsEqual(x, y string) bool {
	if strings.Compare(x, y) != 0 {
		return false
	}
	return true
}
