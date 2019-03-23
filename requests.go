package calendly

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const contentTypeUrlEncoded = "application/x-www-form-urlencoded"
const defaultAPIURL = "https://calendly.com/api/v1"

type Api struct {
	AuthToken string
	// e.g. https://calendly.com/api/v1
	BaseUrl string
}

func New(authToken string) *Api {
	return &Api{AuthToken:authToken, BaseUrl: defaultAPIURL}
}

func (a *Api) Echo() (*EchoResponse, error) {
	var resp EchoResponse
	return &resp, a.request("GET", "/echo", nil, "", &resp)
}

func (a *Api) CreateHook(input CreateHookInput) (*CreateHookResponse, error) {
	data := url.Values{}
	data.Set("url", input.Url)

	for _, ev := range input.Events {
		data.Add("events[]", ev)
	}

	var resp CreateHookResponse
	return &resp, a.request("POST", "/hooks", strings.NewReader(data.Encode()), contentTypeUrlEncoded, &resp)
}

func (a *Api) request(method string, path string, body io.Reader, contentType string, out interface{}) error {
	req, err := http.NewRequest(method, a.BaseUrl+path, body)
	if err != nil {
		return err
	}

	req.Header.Set("X-TOKEN", a.AuthToken)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
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
