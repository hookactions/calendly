package calendly

import (
	"fmt"
	"net/http"
	"time"
)

type BasicResponse struct {
	Status  int    `json:"status,omitempty"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
	// field_name: error_messages_array
	ValidationErrors map[string]interface{} `json:"errors,omitempty"`
}

func (r *BasicResponse) Err() error {
	if r.Status != 0 && r.Status != 200 {
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

	Id int `json:"id"`
}

type CreateHookInput struct {
	URL    string
	Events []string
}

type GetHookInput struct {
	Id int
}

type hookAttributes struct {
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	State     string    `json:"state"`
	Events    []string  `json:"events"`
}

type hook struct {
	Id         int            `json:"id"`
	Type       string         `json:"type"`
	Attributes hookAttributes `json:"attributes"`
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
	Id int
}

type DeleteHookResponse struct {
	BasicResponse
}

type userAttributes struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Email    *string `json:"email"`
	URL      string  `json:"url"`
	Timezone string  `json:"timezone"`
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

type GetEventTypesInput struct {
	Include string
}

type eventTypeAttributes struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Slug        string    `json:"slug"`
	Color       string    `json:"color"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	URL         string    `json:"url"`
}

type eventTypeRelationship struct {
	Owner struct {
		Data struct {
			Type string `json:"type"`
			Id   string `json:"id"`
		} `json:"data"`
	} `json:"owner"`
}

type EventTypesResponse struct {
	BasicResponse

	Data []struct {
		Id            string                 `json:"id"`
		Type          string                 `json:"type"`
		Attributes    eventTypeAttributes    `json:"attributes"`
		Relationships *eventTypeRelationship `json:"relationships,omitempty"`
	} `json:"data"`

	Included []struct {
		Id         string         `json:"id"`
		Type       string         `json:"type"`
		Attributes userAttributes `json:"attributes"`
	} `json:"included,omitempty"`
}
