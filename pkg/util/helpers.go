package util

import "go-deploy/dto/v2/body"

func Float64Pointer(f float64) *float64 {
	return &f
}

func IntPointer(i int) *int {
	return &i
}

func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func GetNames(apiKeys []body.ApiKey) []string {
	var names []string
	for _, apiKey := range apiKeys {
		names = append(names, apiKey.Name)
	}
	return names
}
