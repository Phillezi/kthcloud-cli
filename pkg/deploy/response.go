package deploy

import (
	"fmt"
	"reflect"
)

// ResponseError represents a structured error for API responses.
type ResponseError struct {
	Resource string
	Stage    string
	Code     int
	BaseErr  error
	Message  string
}

func (e *ResponseError) Error() string {
	if e.Code > 0 {
		return fmt.Sprintf("[%s/%s] (%d) %v: %s", e.Resource, e.Stage, e.Code, e.BaseErr, e.Message)
	}
	return fmt.Sprintf("[%s/%s] %v: %s", e.Resource, e.Stage, e.BaseErr, e.Message)
}

func (e *ResponseError) Unwrap() error {
	return e.BaseErr
}

// HandleAPIResponse handles known types explicitly, and falls back to reflection for others.
func HandleAPIResponse(resourceName, stage string, resp any) (any, error) {
	if resp == nil {
		return nil, ErrNilResponse
	}
	switch r := resp.(type) {

	case *GetV2DeploymentsResponse:
		switch {
		case r.JSON200 != nil:
			return r.JSON200, nil

		case r.JSON400 != nil:
			return nil, newResponseError(resourceName, stage, r.StatusCode(), ErrBadRequest,
				fmt.Sprintf("%v", r.JSON400.Errors))

		case r.JSON500 != nil:
			return nil, newResponseError(resourceName, stage, r.StatusCode(), ErrServerError,
				fmt.Sprintf("%v", r.JSON500.Errors))

		default:
			return nil, newResponseError(resourceName, stage, r.StatusCode(), ErrUnexpected,
				string(r.Body))
		}

	default:
		return handleGenericResponse(resourceName, stage, resp)
	}
}

// newResponseError creates a consistent ResponseError.
func newResponseError(resource, stage string, code int, baseErr error, msg string) *ResponseError {
	return &ResponseError{
		Resource: resource,
		Stage:    stage,
		Code:     code,
		BaseErr:  baseErr,
		Message:  msg,
	}
}

func handleGenericResponse(resourceName, stage string, resp any) (any, error) {
	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, newResponseError(resourceName, stage, 0, ErrInvalidResponse,
			fmt.Sprintf("type: %T", resp))
	}

	// Success responses
	for _, name := range []string{"JSON200", "JSON201", "JSON202"} {
		if f := v.FieldByName(name); f.IsValid() && !f.IsNil() {
			return f.Interface(), nil
		}
	}

	// Error responses
	for _, name := range []string{"JSON400", "JSON401", "JSON403", "JSON404", "JSON500"} {
		if f := v.FieldByName(name); f.IsValid() && !f.IsNil() {
			msg := extractErrorMessage(f.Interface())

			var baseErr error
			switch name {
			case "JSON400":
				baseErr = ErrBadRequest
			case "JSON401":
				baseErr = ErrUnauthorized
			case "JSON403":
				baseErr = ErrForbidden
			case "JSON404":
				baseErr = ErrNotFound
			case "JSON500":
				baseErr = ErrServerError
			default:
				baseErr = ErrUnexpected
			}

			code := 0
			if m := v.MethodByName("StatusCode"); m.IsValid() {
				if out := m.Call(nil); len(out) > 0 {
					if c, ok := out[0].Interface().(int); ok {
						code = c
					}
				}
			}

			return nil, newResponseError(resourceName, stage, code, baseErr, msg)
		}
	}

	// Fallback: unknown response type
	statusCode := 0
	if m := v.MethodByName("StatusCode"); m.IsValid() {
		if out := m.Call(nil); len(out) > 0 {
			if c, ok := out[0].Interface().(int); ok {
				statusCode = c
			}
		}
	}
	body := ""
	if f := v.FieldByName("Body"); f.IsValid() && f.Kind() == reflect.Slice {
		body = string(f.Bytes())
	}

	return nil, newResponseError(resourceName, stage, statusCode, ErrUnexpected,
		fmt.Sprintf("status %d: %s", statusCode, body))
}

func extractErrorMessage(errObj any) string {
	v := reflect.ValueOf(errObj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.IsValid() {
		return "<no error info>"
	}

	if f := v.FieldByName("Message"); f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	if f := v.FieldByName("Errors"); f.IsValid() {
		return fmt.Sprintf("%v", f.Interface())
	}
	return fmt.Sprintf("%v", errObj)
}
