package contentful

import (
	"fmt"
	"net/http"
	"net/url"
)

// AppInstallationsService service
type AppInstallationsService service

// AppInstallation model
type AppInstallation struct {
	Sys        *Sys              `json:"sys"`
	Parameters map[string]string `json:"parameters"`
}

// GetVersion returns entity version
func (appInstallation *AppInstallation) GetVersion() int {
	version := 1
	if appInstallation.Sys != nil {
		version = appInstallation.Sys.Version
	}

	return version
}

// List returns an app installations collection
func (service *AppInstallationsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations", spaceID, service.c.Environment)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single app installation
func (service *AppInstallationsService) Get(spaceID, appDefinitionID string) (*AppInstallation, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations/%s", spaceID, service.c.Environment, appDefinitionID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &AppInstallation{}, err
	}

	var installation AppInstallation
	if ok := service.c.do(req, &installation); ok != nil {
		return nil, err
	}

	return &installation, err
}
