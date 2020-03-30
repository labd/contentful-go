package contentful

import (
	"fmt"
	"net/url"
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
	SRC        string      `json:"src"`
	Name       string      `json:"name"`
	FieldTypes []FieldType `json:"fieldTypes"`
	Sidebar    bool        `json:"sidebar"`
}

// FieldType model
type FieldType struct {
	Type string `json:"type"`
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
