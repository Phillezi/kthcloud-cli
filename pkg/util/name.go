package util

import (
	"fmt"
	"regexp"
	"strings"
)

func IsValidDeplName(name string) bool {
	if len(name) > 30 {
		return false
	}

	// Check if the name contains only lowercase alphanumeric characters or '-'
	// and doesn't start or end with a '-'
	re := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	return re.MatchString(name)
}

func GetNameFromUser() (string, error) {
	var name string
	for {
		fmt.Print("Enter a name (less than 30 characters, lowercase, alphanumeric, and - only): ")
		_, err := fmt.Scanln(&name)
		if err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}

		name = strings.TrimSpace(name)

		if len(name) == 0 {
			return "", fmt.Errorf("name cannot be empty")
		}

		if IsValidDeplName(name) {
			return name, nil
		}

		fmt.Println("Invalid name. Please ensure the name follows the RFC 1035 rules (lowercase, alphanumeric, hyphens, and less than 30 characters).")
	}
}
