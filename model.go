package calendly

import (
	"fmt"
	"net/http"
	"time"
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

type GetHookInput struct {
	Id string
}

type hook struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type GetHookResponse struct {
	BasicResponse

	// todo: the documentation shows this as a list, with a single object, is that correct?
	Data []hook `json:"data"`
}

type GetHooksResponse struct {
	BasicResponse

	Data []hook `json:"data"`
}

type DeleteHookInput struct {
	Id string
}

type DeleteHookResponse struct {
	BasicResponse
}

type userAttributes struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Email    string `json:"email"`
	URL      string `json:"url"`
	Timezone string `json:"timezone"`
	Avatar   struct {
		URL string `json:"url"`
	} `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MeResponse struct {
	BasicResponse

	Data struct {
		Id         string         `json:"id"`
		Type       string         `json:"type"`
		Attributes userAttributes `json:"attributes"`
	} `json:"data"`
}
