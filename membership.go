package contentful

import (
	"fmt"
	"net/http"
)

// EntriesService service
type MembershipsService service

//Entry model
type Membership struct {
	locale string
	Sys    *Sys `json:"sys"`
	Fields map[string]interface{}
}

// GetVersion returns entity version
func (membership *Membership) GetVersion() int {
	version := 1
	if membership.Sys != nil {
		version = membership.Sys.Version
	}

	return version
}

// List returns membership collection
func (service *MembershipsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/memberships", spaceID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}
