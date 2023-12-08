package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AppDefinitionsService service
type AppDefinitionsService service

// AppDefinition model
type AppDefinition struct {
	Sys       *Sys        `json:"sys"`
	Name      string      `json:"name"`
	SRC       *string     `json:"src,omitempty"`
	Bundle    *Bundle     `json:"bundle,omitempty"`
	Locations []Locations `json:"locations"`
}

type Bundle struct {
	Sys *Sys `json:"sys"`
}

// Locations model
type Locations struct {
	Location       string          `json:"location"`
	FieldTypes     []FieldType     `json:"fieldTypes,omitempty"`
	NavigationItem *NavigationItem `json:"navigationItem,omitempty"`
}

type FieldType struct {
	Type     string  `json:"type"`
	LinkType *string `json:"linkType,omitempty"`
	Items    *Items  `json:"items,omitempty"`
}

type NavigationItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Items struct {
	Type     string  `json:"type"`
	LinkType *string `json:"linkType,omitempty"`
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

// Get returns a single app definition
func (service *AppDefinitionsService) Get(organizationID, appDefinitionID string) (*AppDefinition, error) {
	path := fmt.Sprintf("/organizations/%s/app_definitions/%s", organizationID, appDefinitionID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return nil, err
	}

	var definition AppDefinition
	if ok := service.c.do(req, &definition); ok != nil {
		return nil, ok
	}

	return &definition, err
}

// Upsert updates or creates a new app definition
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
