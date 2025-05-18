package convert

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/compose-spec/compose-go/v2/types"
)

func HashServices(services types.Services) string {
	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	return Hash(keys...)
}

func Hash(keys ...string) string {
	sort.Strings(keys)
	combinedKeys := strings.Join(keys, "")

	hash := sha256.New()
	hash.Write([]byte(combinedKeys))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
