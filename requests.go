package calendly

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const contentTypeURLEncoded = "application/x-www-form-urlencoded"
const defaultAPIURL = "https://calendly.com/api/v1"

type Api struct {
	AuthToken string
	// e.g. https://calendly.com/api/v1
	BaseURL string
}

func New(authToken string) *Api {
	return &Api{AuthToken: authToken, BaseURL: defaultAPIURL}
}

func (a *Api) Echo() (*EchoResponse, error) {
	var resp EchoResponse
	return &resp, a.request("GET", "/echo", nil, &resp)
}

func (a *Api) CreateHook(input CreateHookInput) (*CreateHookResponse, error) {
	data := url.Values{}
	data.Set("url", input.URL)

	for _, ev := range input.Events {
		data.Add("events[]", ev)
	}

	var resp CreateHookResponse
	return &resp, a.request("POST", "/hooks", &requestOptions{Body: strings.NewReader(data.Encode()), ContentType: contentTypeURLEncoded}, &resp)
}

func (a *Api) GetHook(input GetHookInput) (*GetHookResponse, error) {
	var resp GetHookResponse
	return &resp, a.request("GET", fmt.Sprintf("/hooks/%s", input.Id), nil, &resp)
}

func (a *Api) GetHooks() (*GetHooksResponse, error) {
	var resp GetHooksResponse
	return &resp, a.request("GET", "/hooks", nil, &resp)
}

func (a *Api) DeleteHook(input DeleteHookInput) (*DeleteHookResponse, error) {
	var resp DeleteHookResponse
	return &resp, a.request("DELETE", fmt.Sprintf("/hooks/%d", input.Id), nil, &resp)
}

func (a *Api) Me() (*MeResponse, error) {
	var resp MeResponse
	return &resp, a.request("GET", "/users/me", nil, &resp)
}

func (a *Api) GetEventTypes(input *GetEventTypesInput) (*EventTypesResponse, error) {
	var params map[string]string

	if input != nil {
		params = map[string]string{"include": input.Include}
	}

	var resp EventTypesResponse
	return &resp, a.request("GET", "/users/me/event_types", &requestOptions{QueryParams: params}, &resp)
}

type requestOptions struct {
	Body        io.Reader
	ContentType string
	QueryParams map[string]string
	Headers     map[string]string
}

func (o *requestOptions) Request(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, o.Body)
	if err != nil {
		return nil, err
	}

	if o.Headers != nil {
		for k, v := range o.Headers {
			req.Header.Set(k, v)
		}
	}

	if o.ContentType != "" {
		req.Header.Set("Content-Type", o.ContentType)
	}

	if o.QueryParams != nil {
		q := req.URL.Query()

		for k, v := range o.QueryParams {
			q.Set(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (a *Api) request(method string, path string, opts *requestOptions, out interface{}) error {
	if opts == nil {
		opts = &requestOptions{}
	}

	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}
	opts.Headers["X-TOKEN"] = a.AuthToken

	req, err := opts.Request(method, a.BaseURL+path)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}
