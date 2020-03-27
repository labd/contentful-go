package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// Get returns a single entry task
func (service *EntryTasksService) Get(spaceID, entryID, entryTaskID string) (*EntryTask, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/tasks/%s", spaceID, service.c.Environment, entryID, entryTaskID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &EntryTask{}, err
	}

	var entryTask EntryTask
	if ok := service.c.do(req, &entryTask); ok != nil {
		return nil, err
	}

	return &entryTask, err
}

// Delete the entry task
func (service *EntryTasksService) Delete(spaceID, entryID, entryTaskID string) error {
	path := fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/tasks/%s", spaceID, service.c.Environment, entryID, entryTaskID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}

// Upsert updates or creates a new entry task
func (service *EntryTasksService) Upsert(spaceID, entryID string, entryTask *EntryTask) error {
	bytesArray, err := json.Marshal(entryTask)
	if err != nil {
		return err
	}

	var path string
	var method string

	if entryTask.Sys != nil && entryTask.Sys.CreatedAt != "" {
		path = fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/tasks/%s", spaceID, service.c.Environment, entryID, entryTask.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/environments/%s/entries/%s/tasks", spaceID, service.c.Environment, entryID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(entryTask.GetVersion()))

	return service.c.do(req, entryTask)
}
