package utils

import (
	"encoding/base64"
	"strings"
)

// A function type that reads a byte array used to mock the random reader
type RandomReader func([]byte) (int, error)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int, randomReader RandomReader) ([]byte, error) {
	b := make([]byte, n)
	_, err := randomReader(b)
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
func GenerateRandomString(s int, randomReader RandomReader) (string, error) {
	b, err := GenerateRandomBytes(s, randomReader)
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

// ContainsSubArray checks if a reference array of strings
// contains a sub array of strings.
func ContainsSubArray(sub []string, reference []string) bool {
	exists := make([]bool, len(reference))
	for i := 0; i < len(sub); i++ {
		for j := 0; j < len(reference); j++ {
			if sub[i] == reference[j] {
				exists[j] = true
			}
		}
	}
	for i := 0; i < len(exists); i++ {
		if !exists[i] {
			return false
		}
	}
	return true
}
