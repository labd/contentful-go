package contentful

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
)

// ResourcesService service
type ResourcesService service

// Resource model
type Resource struct {
	Sys *Sys `json:"sys"`
}

// Get returns a single resource/upload
func (service *ResourcesService) Get(spaceID, resourceID string) (*Resource, error) {
	path := fmt.Sprintf("/spaces/%s/uploads/%s", spaceID, resourceID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Resource{}, err
	}

	var resource Resource
	if ok := service.c.do(req, &resource); ok != nil {
		return nil, err
	}

	return &resource, err
}

// Create creates an upload resource
func (service *ResourcesService) Create(spaceID, filePath string) error {
	bytesArray, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/spaces/%s/uploads", spaceID)
	method := "POST"

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	return service.c.do(req, bytesArray)
}
