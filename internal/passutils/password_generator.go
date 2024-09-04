package passutils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GeneratePassword(length int) (string, error) {
	if length < 6 {
		return "", fmt.Errorf("minimal password length is 6")
	}

	lowercase := "abcdefghijklmnopqrstuvwxyz"
	uppercase := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers := "0123456789"
	special := "~!@#$%^&*()_-+={[}}|:;<,>.?/"
	all := lowercase + uppercase + numbers + special

	var password strings.Builder
	password.Grow(length)

	for _, charset := range []string{lowercase, uppercase, numbers, special} {
		if err := addRandomChar(&password, charset); err != nil {
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
	}

	for i := 4; i < length; i++ {
		if err := addRandomChar(&password, all); err != nil {
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
	}

	shuffled := shuffleString(password.String())

	return shuffled, nil
}

func addRandomChar(builder *strings.Builder, charset string) error {
	charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return err
	}

	builder.WriteByte(charset[charIndex.Int64()])
	return nil
}

func shuffleString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		runes[i], runes[j.Int64()] = runes[j.Int64()], runes[i]
	}
	return string(runes)
}
