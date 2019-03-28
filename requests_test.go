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
	eventTypesDataWithOwner = `{
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
      },
      "relationships":{
        "owner":{
          "data":{
            "type":"users",
            "id":"XXXXXXXXXXXXXXXX"
          }
        }
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
      },
      "relationships":{
        "owner":{
          "data":{
            "type":"users",
            "id":"XXXXXXXXXXXXXXXX"
          }
        }
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
      },
      "relationships":{
        "owner":{
          "data":{
            "type":"teams",
            "id":"ZZZZZZZZZZZZZZZZ"
          }
        }
      }
    }
  ],
  "included":[
    {
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
    },
    {
      "type":"teams",
      "id":"ZZZZZZZZZZZZZZZZ",
      "attributes":{
        "name":"ACME Team",
        "slug":"acme-team",
        "email":null,
        "url":"https://calendly.com/acme-team",
        "timezone":"America/New_York",
        "avatar":{
          "url":"https://d3v0px0pttie1i.cloudfront.net/uploads/team/avatar/2682/9e56907a.gif"
        },
        "created_at":"2016-03-24T16:09:01Z",
        "updated_at":"2016-03-24T16:09:01Z"
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
	hookDetailData = `{
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
    }
  ]
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
		case "/hooks/123":
			if r.Method == "DELETE" {
				_, err := w.Write([]byte(deleteHookData))
				assert.NoError(t, err)
				return
			} else if r.Method == "GET" {
				_, err := w.Write([]byte(hookDetailData))
				assert.NoError(t, err)
				return
			}
		case "/users/me":
			if r.Method == "GET" {
				_, err := w.Write([]byte(meData))
				assert.NoError(t, err)
				return
			}
		case "/users/me/event_types":
			if r.Method == "GET" {
				include := r.URL.Query().Get("include")

				if include != "" {
					_, err := w.Write([]byte(eventTypesDataWithOwner))
					assert.NoError(t, err)
					return
				}

				_, err := w.Write([]byte(eventTypesData))
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

	assert.Equal(t, hook{
		Type: "hooks",
		Id:   12345,
		Attributes: hookAttributes{
			URL:       "http://foo.bar/1",
			CreatedAt: time.Date(2016, 8, 23, 19, 15, 24, 0, time.UTC),
			State:     "active",
			Events:    []string{"invitee.created", "invitee.canceled"},
		},
	}, resp.Data[0])
	assert.Equal(t, hook{
		Type: "hooks",
		Id:   1234,
		Attributes: hookAttributes{
			URL:       "http://foo.bar/2",
			CreatedAt: time.Date(2016, 2, 11, 19, 10, 12, 0, time.UTC),
			State:     "disabled",
			Events:    []string{"invitee.created"},
		},
	}, resp.Data[1])
}

func TestApi_DeleteHook(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}

	resp, err := api.DeleteHook(DeleteHookInput{Id: 123})
	require.NoError(t, err)

	assert.NoError(t, resp.Err())
}

func TestApi_Me(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}

	resp, err := api.Me()
	require.NoError(t, err)

	email := "janedoe30305@gmail.com"

	expected := &MeResponse{}
	expected.Data.Id = "XXXXXXXXXXXXXXXX"
	expected.Data.Type = "users"
	expected.Data.Attributes = userAttributes{
		Name:      "Jane Doe",
		Slug:      "janedoe",
		Email:     &email,
		URL:       "https://calendly.com/janedoe",
		Timezone:  "America/New_York",
		CreatedAt: time.Date(2015, 6, 16, 18, 46, 52, 0, time.UTC),
		UpdatedAt: time.Date(2016, 8, 23, 19, 40, 7, 0, time.UTC),
	}
	expected.Data.Attributes.Avatar.URL = "https://d3v0px0pttie1i.cloudfront.net/uploads/user/avatar/68272/78fb9f5e.jpg"

	assert.Equal(t, expected, resp)
}

func TestApi_GetEventTypes(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}

	resp, err := api.GetEventTypes(nil)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 3)
	assert.Empty(t, resp.Included)

	assert.Equal(t, "event_types", resp.Data[0].Type)
	assert.Equal(t, "AAAAAAAAAAAAAAAA", resp.Data[0].Id)
	assert.Equal(t, eventTypeAttributes{
		Name:        "15 Minute Meeting",
		Description: "",
		Duration:    30,
		Slug:        "15min",
		Color:       "#fff200",
		Active:      true,
		CreatedAt:   time.Date(2015, 6, 16, 18, 46, 53, 0, time.UTC),
		UpdatedAt:   time.Date(2016, 8, 23, 19, 27, 52, 0, time.UTC),
		URL:         "https://calendly.com/janedoe/15min",
	}, resp.Data[0].Attributes)

	t.Run("Include", func(t *testing.T) {
		resp, err = api.GetEventTypes(&GetEventTypesInput{Include: "owner"})
		require.NoError(t, err)

		assert.Len(t, resp.Data, 3)
		assert.Len(t, resp.Included, 2)

		assert.Equal(t, "event_types", resp.Data[0].Type)
		assert.Equal(t, "AAAAAAAAAAAAAAAA", resp.Data[0].Id)
		assert.Equal(t, eventTypeAttributes{
			Name:        "15 Minute Meeting",
			Description: "",
			Duration:    30,
			Slug:        "15min",
			Color:       "#fff200",
			Active:      true,
			CreatedAt:   time.Date(2015, 6, 16, 18, 46, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2016, 8, 23, 19, 27, 52, 0, time.UTC),
			URL:         "https://calendly.com/janedoe/15min",
		}, resp.Data[0].Attributes)

		expected := &eventTypeRelationship{}
		expected.Owner.Data.Type = "users"
		expected.Owner.Data.Id = "XXXXXXXXXXXXXXXX"

		assert.Equal(t, expected, resp.Data[0].Relationships)

		assert.Equal(t, "teams", resp.Included[1].Type)
		assert.Equal(t, "ZZZZZZZZZZZZZZZZ", resp.Included[1].Id)

		expectedTeamAttributes := userAttributes{
			Name:      "ACME Team",
			Slug:      "acme-team",
			Email:     nil,
			URL:       "https://calendly.com/acme-team",
			Timezone:  "America/New_York",
			CreatedAt: time.Date(2016, 3, 24, 16, 9, 1, 0, time.UTC),
			UpdatedAt: time.Date(2016, 3, 24, 16, 9, 1, 0, time.UTC),
		}
		expectedTeamAttributes.Avatar.URL = "https://d3v0px0pttie1i.cloudfront.net/uploads/team/avatar/2682/9e56907a.gif"

		assert.Equal(t, expectedTeamAttributes, resp.Included[1].Attributes)
	})
}

func TestApi_GetHook(t *testing.T) {
	server := testApiServer(t, "123")
	defer server.Close()

	api := &Api{BaseURL: server.URL, AuthToken: "123"}
	resp, err := api.GetHook(GetHookInput{Id: 123})
	require.NoError(t, err)

	assert.Equal(t, hook{
		Id:   12345,
		Type: "hooks",
		Attributes: hookAttributes{
			URL:       "http://foo.bar/1",
			CreatedAt: time.Date(2016, 8, 23, 19, 15, 24, 0, time.UTC),
			State:     "active",
			Events:    []string{"invitee.created", "invitee.canceled"},
		},
	}, resp.Data[0])
}
