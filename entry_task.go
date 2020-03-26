package contentful

import (
	"fmt"
	"net/http"
)

// EntryTasksService service
type EntryTasksService service

// EntryTask model
type EntryTask struct {
	Sys        *Sys       `json:"sys"`
	Body       string     `json:"body"`
	Status     string     `json:"status"`
	AssignedTo AssignedTo `json:"assignedTo"`
}

// AssignedTo model
type AssignedTo struct {
	Sys Sys `json:"sys"`
}

// GetVersion returns entity version
func (entryTask *EntryTask) GetVersion() int {
	version := 1
	if entryTask.Sys != nil {
		version = entryTask.Sys.Version
	}

	return version
}

// List returns entry tasks collection
func (service *EntryTasksService) List(spaceID, entryID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/tasks", spaceID, service.c.Environment, entryID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}
