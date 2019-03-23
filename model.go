package calendly

import (
	"fmt"
	"net/http"
)

type Attributes map[string]interface{}

type BasicResponse struct {
	Status  int    `json:"status,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
	// field_name: error_messages_array
	ValidationErrors map[string]interface{} `json:"errors,omitempty"`
}

func (r *BasicResponse) Err() error {
	if r.Status != 200 {
		return fmt.Errorf(r.Message)
	}
	return nil
}

func (r *BasicResponse) Unauthorized() bool {
	return r.Status == http.StatusUnauthorized
}

func (r *BasicResponse) Forbidden() bool {
	return r.Status == http.StatusForbidden
}

func (r *BasicResponse) HasValidationError() bool {
	return r.Status == http.StatusUnprocessableEntity
}

type EchoResponse struct {
	BasicResponse

	Email string `json:"email"`
}

type CreateHookResponse struct {
	BasicResponse

	Id string `json:"id"`
}

type CreateHookInput struct {
	Url    string
	Events []string
}
