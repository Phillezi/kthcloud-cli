package run

import (
	"strings"

	"math/rand/v2"
)

func GenerateRandomName(minLen, maxLen int) string {
	if minLen < 3 {
		minLen = 3
	}
	if maxLen > 30 {
		maxLen = 30
	}

	letters := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	validChars := letters + digits + "-"

	nameLen := rand.IntN(maxLen-minLen+1) + minLen
	name := strings.Builder{}
	name.Grow(nameLen)

	name.WriteByte(letters[rand.IntN(len(letters))])

	for i := 1; i < nameLen-1; i++ {
		name.WriteByte(validChars[rand.IntN(len(validChars))])
	}

	lastChars := letters + digits
	name.WriteByte(lastChars[rand.IntN(len(lastChars))])

	return name.String()
}
