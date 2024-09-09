package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func HandleResponse(response *resty.Response) error {
	log.Infoln(response)
	// Handle different HTTP status codes
	switch response.StatusCode() {
	case http.StatusOK:

	case http.StatusNotFound:
		return fmt.Errorf("resource not found: %s", response.Status())
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized access: %s", response.Status())
	case http.StatusForbidden:
		return fmt.Errorf("forbidden access: %s", response.Status())
	case http.StatusInternalServerError:
		return fmt.Errorf("server error: %s", response.Status())
	default:
		return fmt.Errorf("unexpected status code: %d", response.StatusCode())
	}

	// Unmarshal the response to a map
	var responseMap map[string]interface{}
	if err := json.Unmarshal([]byte(response.String()), &responseMap); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for errors in the response
	if errors, ok := responseMap["errors"]; ok {
		// If errors are present, format them
		var errorMessages []string
		switch errs := errors.(type) {
		case []interface{}:
			for _, err := range errs {
				if errMap, ok := err.(map[string]interface{}); ok {
					code, _ := errMap["code"].(string)
					msg, _ := errMap["msg"].(string)
					errorMessages = append(errorMessages, fmt.Sprintf("Code: %s, Message: %s", code, msg))
				}
			}
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("Unexpected error format: %v", errors))
		}

		// Log the errors
		errorMsg := strings.Join(errorMessages, "; ")
		log.Errorf("Response contains errors: %s", errorMsg)
		return fmt.Errorf("response contains errors: %s", errorMsg)
	}

	// No errors in the response
	return nil
}
