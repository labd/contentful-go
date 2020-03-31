package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AppDefinitionsService service
type AppDefinitionsService service

// AppDefinition model
type AppDefinition struct {
	Sys       *Sys        `json:"sys"`
	Name      string      `json:"name"`
	SRC       string      `json:"src"`
	Locations []Locations `json:"locations"`
}

// Locations model
type Locations struct {
	Location string `json:"location"`
}

// GetVersion returns entity version
func (appDefinition *AppDefinition) GetVersion() int {
	version := 1
	if appDefinition.Sys != nil {
		version = appDefinition.Sys.Version
	}

	return version
}

// List returns an app definitions collection
func (service *AppDefinitionsService) List(organizationID string) *Collection {
	path := fmt.Sprintf("/organizations/%s/app_definitions", organizationID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single entry
func (service *AppDefinitionsService) Get(organizationID, appDefinitionID string) (*AppDefinition, error) {
	path := fmt.Sprintf("/organizations/%s/app_definitions/%s", organizationID, appDefinitionID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &AppDefinition{}, err
	}

	var definition AppDefinition
	if ok := service.c.do(req, &definition); ok != nil {
		return nil, err
	}

	return &definition, err
}

// Upsert updates or creates a new entry
func (service *AppDefinitionsService) Upsert(organizationID string, definition *AppDefinition) error {
	bytesArray, err := json.Marshal(definition)
	if err != nil {
		return err
	}

	var path string
	var method string

	if definition.Sys != nil && definition.Sys.ID != "" {
		path = fmt.Sprintf("/organizations/%s/app_definitions/%s", organizationID, definition.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/organizations/%s/app_definitions", organizationID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(definition.GetVersion()))

	return service.c.do(req, definition)
}

// Delete the app definition
func (service *AppDefinitionsService) Delete(organizationID, appDefinitionID string) error {
	path := fmt.Sprintf("/organizations/%s/app_definitions/%s", organizationID, appDefinitionID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
