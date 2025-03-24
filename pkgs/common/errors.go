package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/labd/contentful-go/pkgs/model"
)

// ErrorResponse model
type ErrorResponse struct {
	Sys       *model.BaseSys `json:"sys"`
	Message   string         `json:"message,omitempty"`
	RequestID string         `json:"requestId,omitempty"`
	Details   *ErrorDetails  `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// ErrorDetails model
type ErrorDetails struct {
	Errors  []*ErrorDetail `json:"errors,omitempty"`
	Reasons string         `json:"reasons,omitempty"`
}

func (e *ErrorDetails) UnmarshalJSON(data []byte) error {
	var unmarshaled interface{}

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	switch reflect.ValueOf(unmarshaled).Kind() {
	case reflect.String:
		e.Errors = append(e.Errors, &ErrorDetail{
			Details: unmarshaled.(string),
		})

	case reflect.Map:
		if _, ok := unmarshaled.(map[string]any)["errors"]; ok {
			intermedidateStruct := &struct {
				Errors []*ErrorDetail `json:"errors,omitempty"`
			}{}

			err = json.Unmarshal(data, &intermedidateStruct)
			if err != nil {
				return err
			}

			e.Errors = intermedidateStruct.Errors

		}

		if val, ok := unmarshaled.(map[string]any)["reasons"]; ok {
			e.Reasons = val.(string)

		}
	}

	return nil
}

// ErrorDetail model
type ErrorDetail struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	Path        interface{} `json:"path,omitempty"`
	Details     string      `json:"details,omitempty"`
	Value       interface{} `json:"value,omitempty"`
	Conflicting []*struct {
		Sys model.BaseSys `json:"sys,omitempty"`
	} `json:"conflicting,omitempty"`
}

// APIError model
type APIError struct {
	req *http.Request
	res *http.Response
	Err *ErrorResponse
}

func NewApiError(req *http.Request, res *http.Response, err *ErrorResponse) APIError {
	return APIError{
		req: req,
		res: res,
		Err: err,
	}
}

// AccessTokenInvalidError for 401 errors
type AccessTokenInvalidError struct {
	APIError
}

func (e AccessTokenInvalidError) Error() string {
	return e.APIError.Err.Message
}

// VersionMismatchError for 409 errors
type VersionMismatchError struct {
	APIError
}

func (e VersionMismatchError) Error() string {
	return "Version " + e.APIError.req.Header.Get("X-Contentful-Version") + " is mismatched"
}

// ValidationFailedError model
type ValidationFailedError struct {
	APIError
}

func (e ValidationFailedError) Error() string {
	msg := bytes.Buffer{}

	for _, err := range e.APIError.Err.Details.Errors {
		if path, ok := getPathAsString(err.Path); ok {
			msg.WriteString(fmt.Sprintf("Value \"%s\" in path \"%s\" with details: \"%s\"\n", err.Value, *path, err.Details))
			continue
		}

		msg.WriteString(fmt.Sprintf("Value %s in path %+v %s\n", err.Value, err.Path, err.Details))
	}

	return msg.String()
}

func getPathAsString(path any) (*string, bool) {
	switch x := path.(type) {
	case []string:
		res := strings.Join(x, ".")
		return &res, true
	case []any:
		var res []string

		for _, val := range x {
			switch val.(type) {
			case string:
				res = append(res, val.(string))

			default:
				res = append(res, fmt.Sprintf("%+v", val))
			}
		}

		ret := strings.Join(res, ".")
		return &ret, true
	default:
		return nil, false
	}
}

// NotFoundError for 404 errors
type NotFoundError struct {
	APIError
}

func (e NotFoundError) Error() string {
	return "the requested resource can not be found"
}

// RateLimitExceededError for rate limit errors
type RateLimitExceededError struct {
	APIError
}

func (e RateLimitExceededError) Error() string {
	return e.APIError.Err.Message
}

// BadRequestError error model for bad request responses
type BadRequestError struct{}

// InvalidQueryError error model for invalid query responses
type InvalidQueryError struct{}

// AccessDeniedError error model for access denied responses
type AccessDeniedError struct{}

// ServerError error model for server error responses
type ServerError struct{}

// InvalidEntryError model
type InvalidEntryError struct {
	APIError
}

func (e InvalidEntryError) Error() string {
	msg := bytes.Buffer{}

	for _, err := range e.APIError.Err.Details.Errors {
		msg.WriteString(fmt.Sprintf("%s\n", err.Details))
	}

	return msg.String()
}
