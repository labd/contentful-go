package contentful

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
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
		return nil, ok
	}

	return &resource, err
}

// Create creates an upload resource
func (service *ResourcesService) Create(spaceID, filePath string) error {
	bytesArray, err := os.ReadFile(filePath)
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

// Delete the resource
func (service *ResourcesService) Delete(spaceID, resourceID string) error {
	path := fmt.Sprintf("/spaces/%s/uploads/%s", spaceID, resourceID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
