package contentful

import (
	"fmt"
	"net/http"
)

// ScheduledActionsService service
type ScheduledActionsService service

// ScheduledAction model
type ScheduledAction struct {
	Sys          *Sys              `json:"sys"`
	Entity       Entity            `json:"entity"`
	Environment  Environment       `json:"environment"`
	ScheduledFor map[string]string `json:"scheduledFor"`
	Action       string            `json:"action"`
}

// Entity model
type Entity struct {
	Sys *Sys `json:"sys"`
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
