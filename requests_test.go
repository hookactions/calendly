package calendly

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	echoData         = `{"email": "test@example.com"}`
	createHookData   = `{"id": "123"}`
	invalidTokenData = `{"status": 401, "type": "authentication_error", "message": "Invalid token"}`
)

func testApiServer(t *testing.T, apiToken string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-TOKEN")
		if token == "" || token != apiToken {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(invalidTokenData))
			assert.NoError(t, err)
			return
		}
		switch r.URL.Path {
		case "/echo":
			assert.Equal(t, "GET", r.Method)
			_, err := w.Write([]byte(echoData))
			assert.NoError(t, err)
			return
		case "/hooks":
			assert.Equal(t, "POST", r.Method)
			defer r.Body.Close()

			assert.NoError(t, r.ParseForm())

			events, ok := r.PostForm["events[]"]
			assert.True(t, ok)
			assert.True(t, len(events) > 0)

			_, ok = r.PostForm["url"]
			assert.True(t, ok)

			_, err := w.Write([]byte(createHookData))
			assert.NoError(t, err)
			return
		}
		http.NotFound(w, r)
	}))
}

func TestBasicResponse_Err(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "foo"}
	resp, err := api.Echo()
	require.NoError(t, err)

	assert.True(t, resp.Unauthorized())
	assert.Equal(t, "authentication_error", resp.Type)
	assert.Equal(t, "Invalid token", resp.Message)
}

func TestApi_Echo(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}
	resp, err := api.Echo()
	require.NoError(t, err)

	assert.Equal(t, "test@example.com", resp.Email)
}

func TestApi_CreateHook(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}
	input := CreateHookInput{URL: "https://example.com", Events: []string{"invitee.created", "invitee.canceled"}}
	resp, err := api.CreateHook(input)
	require.NoError(t, err)

	assert.Equal(t, "123", resp.Id)
}
