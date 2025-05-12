package utils

import (
	"math/rand"
	"strings"
)

// GenerateAccountNumber generates a random account number with a given length. The generated number is a string whose value "may" be unique
func GenerateAccountNumber(length int) string {
	domain := "0123456789"

	// ensure that each random generation results in a new value
	var numStr strings.Builder

	for range length {
		val := domain[rand.Intn(len(domain))]
		numStr.WriteString(string(val))
	}
	return numStr.String()
}
