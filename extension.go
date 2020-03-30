package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ExtensionsService service
type ExtensionsService service

// Extension model
type Extension struct {
	Sys       *Sys             `json:"sys"`
	Extension ExtensionDetails `json:"extension"`
}

// ExtensionDetails model
type ExtensionDetails struct {
	SRC        string       `json:"src"`
	Name       string       `json:"name"`
	FieldTypes []FieldTypes `json:"fieldTypes"`
	Sidebar    bool         `json:"sidebar"`
}

// FieldTypes model
type FieldTypes struct {
	Type string `json:"type"`
}

// GetVersion returns entity version
func (extension *Extension) GetVersion() int {
	version := 1
	if extension.Sys != nil {
		version = extension.Sys.Version
	}

	return version
}

// List returns an extensions collection
func (service *ExtensionsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/extensions", spaceID, service.c.Environment)

	req, err := service.c.newRequest("GET", path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single extension
func (service *ExtensionsService) Get(spaceID, extensionID string) (*Extension, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/extensions/%s", spaceID, service.c.Environment, extensionID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Extension{}, err
	}

	var extension Extension
	if ok := service.c.do(req, &extension); ok != nil {
		return nil, err
	}

	return &extension, err
}

// Upsert updates or creates a new extension
func (service *ExtensionsService) Upsert(spaceID string, e *Extension) error {
	bytesArray, err := json.Marshal(e)
	if err != nil {
		return err
	}

	var path string
	var method string

	if e.Sys != nil && e.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s/environments/%s/extensions/%s", spaceID, service.c.Environment, e.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/environments/%s/extensions", spaceID, service.c.Environment)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(e.GetVersion()))

	return service.c.do(req, e)
}

// Delete the extension
func (service *ExtensionsService) Delete(spaceID string, extensionID string) error {
	path := fmt.Sprintf("/spaces/%s/environments/%s/extensions/%s", spaceID, service.c.Environment, extensionID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
