package contentful

import (
	"fmt"
	"net/http"
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
