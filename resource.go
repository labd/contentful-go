package contentful

import (
	"fmt"
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
