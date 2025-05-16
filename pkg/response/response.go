package response

import (
	"fmt"
	"strings"

	"github.com/Phillezi/kthcloud-cli/pkg/util"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func IsError(response string) error {
	if strings.HasPrefix(response, "{\"errors\":") {
		errors, err := util.ProcessResponse[ErrorResponse](response)
		if err != nil {
			// could not convert to errorresponse but still contains errors
			return err
		}
		return fmt.Errorf("error when trying to create deployment %v", *errors)
	}
	return nil
}
