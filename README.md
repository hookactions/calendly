# calendly

[![CircleCI](https://circleci.com/gh/hookactions/calendly/tree/master.svg?style=svg)](https://circleci.com/gh/hookactions/calendly/tree/master)

API wrapper for Calendly. https://developer.calendly.com/

# Use

_For brevity, error checking is omitted._

```go
package main

import (
	"fmt"
	
	"github.com/hookactions/calendly"
)

func main() {
	api := calendly.New("api-token")
	
	// 1. Echo, test auth is working
	echoResp, _ := api.Echo()
	fmt.Printf("Logged in as: %s\n", echoResp.Email)
	
	// 2. Create a webhook
	createResp, _ := api.CreateHook(calendly.CreateHookInput{
		URL: "https://foo.bar/my-hook",
		Events: []string{"invitee.created", "invitee.canceled"},
	})
	fmt.Printf("Created hook with id: %d\n", createResp.Id)
	
	// 3. Get webhook by id
	getResp, _ := api.GetHook(calendly.GetHookInput{Id: 1})
	fmt.Printf("Got hook, state is: %s\n", getResp.Data[0].Attributes.State)
	
	// 4. Get all webhooks
	listResp, _ := api.GetHooks()
	fmt.Printf("Got %d hooks\n", len(listResp.Data))
	
	// 5. Delete webhook
	api.DeleteHook(calendly.DeleteHookInput{Id: 1})
	fmt.Println("deleted")
	
	// 6. Get user event types
	eventTypes, _ := api.GetEventTypes(nil)
	fmt.Printf("Got %d event types\n", len(eventTypes.Data))
	
	// 7. Get current user data
	me, _ := api.Me()
	fmt.Printf("Current user name: %s", me.Data.Attributes.Name)
}
```

## Error checking

All responses have the following methods for checking for certain errors.

```go
package main

import (
	"fmt"
	
	"github.com/hookactions/calendly"
)


func main() {
	api := calendly.New("api-token")
	createResp, _ := api.CreateHook(calendly.CreateHookInput{
		URL: "https://foo.bar/my-hook", 
		Events: []string{"invitee.created", "invitee.canceled"},
	})
    
    if createResp.Err() != nil {
    	fmt.Printf("Unauthorized? %v\n", createResp.Unauthorized())
    	fmt.Printf("Forbidden? %v\n", createResp.Forbidden())
    	fmt.Printf("Validation error? %v\n", createResp.HasValidationError())
    }
    
    if createResp.HasValidationError() {
    	fmt.Printf("Validation errors: %#v\n", createResp.ValidationErrors)
    }
}
```
