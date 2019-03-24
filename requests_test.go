package calendly

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	echoData         = `{"email": "test@example.com"}`
	createHookData   = `{"id": 123}`
	invalidTokenData = `{"status": 401, "type": "authentication_error", "message": "Invalid token"}`
	hookListData     = `{
  "data":[
    {
      "type":"hooks",
      "id":12345,
      "attributes":{
        "url":"http://foo.bar/1",
        "created_at":"2016-08-23T19:15:24Z",
        "state":"active",
        "events":[
          "invitee.created",
          "invitee.canceled"
        ]
      }
    },
    {
      "type":"hooks",
      "id":1234,
      "attributes":{
        "url":"http://foo.bar/2",
        "created_at":"2016-02-11T19:10:12Z",
        "state":"disabled",
        "events":[
          "invitee.created"
        ]
      }
    }
  ]
}`
	deleteHookData = `null`
	eventTypesData = `{
  "data":[
    {
      "type":"event_types",
      "id":"AAAAAAAAAAAAAAAA",
      "attributes":{
        "name":"15 Minute Meeting",
        "description":"",
        "duration":30,
        "slug":"15min",
        "color":"#fff200",
        "active":true,
        "created_at":"2015-06-16T18:46:53Z",
        "updated_at":"2016-08-23T19:27:52Z",
        "url":"https://calendly.com/janedoe/15min"
      }
    },
    {
      "type":"event_types",
      "id":"BBBBBBBBBBBBBBBB",
      "attributes":{
        "name":"30 Minute Meeting",
        "description":"",
        "duration":30,
        "slug":"30min",
        "color":"#74daed",
        "active":true,
        "created_at":"2015-06-16T18:46:53Z",
        "updated_at":"2016-06-02T16:26:44Z",
        "url":"https://calendly.com/janedoe/30min"
      }
    },
    {
      "type":"event_types",
      "id":"CCCCCCCCCCCCCCCC",
      "attributes":{
        "name":"Sales call",
        "description":null,
        "duration":30,
        "slug":"sales-call",
        "color":"#c5c1ff",
        "active":false,
        "created_at":"2016-06-23T20:13:17Z",
        "updated_at":"2016-06-23T20:13:22Z",
        "url":"https://calendly.com/acme-team/sales-call"
      }
    }
  ]
}`
	meData = `{
  "data":{
    "type":"users",
    "id":"XXXXXXXXXXXXXXXX",
    "attributes":{
      "name":"Jane Doe",
      "slug":"janedoe",
      "email":"janedoe30305@gmail.com",
      "url":"https://calendly.com/janedoe",
      "timezone":"America/New_York",
      "avatar":{
        "url":"https://d3v0px0pttie1i.cloudfront.net/uploads/user/avatar/68272/78fb9f5e.jpg"
      },
      "created_at":"2015-06-16T18:46:52Z",
      "updated_at":"2016-08-23T19:40:07Z"
    }
  }
}`
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
			if r.Method == "POST" {
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
			} else if r.Method == "GET" {
				_, err := w.Write([]byte(hookListData))
				assert.NoError(t, err)
				return
			}
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

	assert.Equal(t, 123, resp.Id)
}

func TestApi_GetHooks(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}

	resp, err := api.GetHooks()
	require.NoError(t, err)

	assert.Len(t, resp.Data, 2)

	assert.Equal(t, resp.Data[0], hook{
		Type: "hooks",
		Id:   12345,
		Attributes: hookAttributes{
			URL:       "http://foo.bar/1",
			CreatedAt: time.Date(2016, 8, 23, 19, 15, 24, 0, time.UTC),
			State:     "active",
			Events:    []string{"invitee.created", "invitee.canceled"},
		},
	})
	assert.Equal(t, resp.Data[1], hook{
		Type: "hooks",
		Id:   1234,
		Attributes: hookAttributes{
			URL:       "http://foo.bar/2",
			CreatedAt: time.Date(2016, 2, 11, 19, 10, 12, 0, time.UTC),
			State:     "disabled",
			Events:    []string{"invitee.created"},
		},
	})
}
