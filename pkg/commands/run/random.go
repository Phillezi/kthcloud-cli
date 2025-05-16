package run

import (
	"strings"
	"time"

	"golang.org/x/exp/rand"
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

	rand.Seed(uint64(time.Now().UnixNano()))

	nameLen := rand.Intn(maxLen-minLen+1) + minLen
	name := strings.Builder{}
	name.Grow(nameLen)

	name.WriteByte(letters[rand.Intn(len(letters))])

	for i := 1; i < nameLen-1; i++ {
		name.WriteByte(validChars[rand.Intn(len(validChars))])
	}

	lastChars := letters + digits
	name.WriteByte(lastChars[rand.Intn(len(lastChars))])

	return name.String()
}
