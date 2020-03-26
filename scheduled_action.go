package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ScheduledActionsService service
type ScheduledActionsService service

// ScheduledAction model
type ScheduledAction struct {
	Sys          *Sys              `json:"sys"`
	Entity       Entity            `json:"entity"`
	Environment  EnvironmentLink   `json:"environment"`
	ScheduledFor map[string]string `json:"scheduledFor"`
	Action       string            `json:"action"`
}

// Entity model
type Entity struct {
	Sys Sys `json:"sys"`
}

// EnvironmentLink model
type EnvironmentLink struct {
	Sys Sys `json:"sys"`
}

// GetVersion returns entity version
func (scheduledAction *ScheduledAction) GetVersion() int {
	version := 1
	if scheduledAction.Sys != nil {
		version = scheduledAction.Sys.Version
	}

	return version
}

// List returns scheduled actions collection
func (service *ScheduledActionsService) List(spaceID, entryID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/scheduled_actions?entity.sys.id=%s&environment.sys.id=%s", spaceID, entryID, service.c.Environment)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Delete the scheduled action
func (service *ScheduledActionsService) Delete(spaceID, entryID, scheduledActionID string) error {
	path := fmt.Sprintf("/spaces/%s/scheduled_actions/%s?entity.sys.id=%s&environment.sys.id=%s", spaceID, scheduledActionID, entryID, service.c.Environment)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Create creates a new scheduled actions
func (service *ScheduledActionsService) Create(spaceID, entryID string, scheduledAction *ScheduledAction) error {
	bytesArray, err := json.Marshal(scheduledAction)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/spaces/%s/scheduled_actions?entity.sys.id=%s&environment.sys.id=%s", spaceID, entryID, service.c.Environment)
	method := "POST"

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(scheduledAction.GetVersion()))

	return service.c.do(req, scheduledAction)
}
