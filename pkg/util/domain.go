package util

import (
	"errors"
	"strings"
)

// GetCommonDomain returns the common domain from multiple domains by removing subdomains.
// It assumes domains provided are hosts without a scheme.
func GetCommonDomain(domains ...string) (string, error) {
	if len(domains) == 0 {
		return "", errors.New("no domains provided")
	}

	commonParts := strings.Split(domains[0], ".")

	for _, domain := range domains[1:] {
		parts := strings.Split(domain, ".")
		commonParts = findCommonSuffix(commonParts, parts)

		if len(commonParts) == 0 {
			return "", errors.New("no common domain found")
		}
	}

	commonDomain := strings.Join(reverse(commonParts), ".")
	return commonDomain, nil
}

// findCommonSuffix finds the longest common suffix between two domain part arrays
func findCommonSuffix(parts1, parts2 []string) []string {
	i, j := len(parts1)-1, len(parts2)-1
	var common []string

	for ; i >= 0 && j >= 0 && parts1[i] == parts2[j]; j-- {
		common = append(common, parts1[i])
		i--
	}

	return common
}

func reverse(arr []string) (reversed []string) {
	for i := len(arr) - 1; i >= 0; i-- {
		reversed = append(reversed, arr[i])
	}
	return reversed
}
