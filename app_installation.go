package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// AppInstallationsService service
type AppInstallationsService service

// AppInstallation model
type AppInstallation struct {
	Sys        *Sys           `json:"sys"`
	Parameters map[string]any `json:"parameters"`
}

// List returns an app installations collection
func (service *AppInstallationsService) List(spaceID string, environment string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations", spaceID, environment)

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
func (service *AppInstallationsService) Get(spaceID, appInstallationID string, environment string) (*AppInstallation, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations/%s", spaceID, environment, appInstallationID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &AppInstallation{}, err
	}

	var installation AppInstallation
	if ok := service.c.do(req, &installation); ok != nil {
		return nil, ok
	}

	return &installation, err
}

// Upsert updates or creates a new app installation
func (service *AppInstallationsService) Upsert(spaceID, appInstallationID string, installation *AppInstallation, environment string) error {
	bytesArray, err := json.Marshal(installation)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations/%s", spaceID, environment, appInstallationID)

	req, err := service.c.newRequest("PUT", path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	return service.c.do(req, installation)
}

// Delete the app installation
func (service *AppInstallationsService) Delete(spaceID, appInstallationID string, environment string) error {
	path := fmt.Sprintf("/spaces/%s/environments/%s/app_installations/%s", spaceID, environment, appInstallationID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
