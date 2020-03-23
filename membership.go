package contentful

import (
	"fmt"
	"net/http"
	"net/url"
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
	path := fmt.Sprintf("/spaces/%s/space_memberships", spaceID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single membership
func (service *MembershipsService) Get(spaceID, membershipID string) (*Membership, error) {
	path := fmt.Sprintf("/spaces/%s/space_memberships/%s", spaceID, membershipID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Membership{}, err
	}

	var membership Membership
	if ok := service.c.do(req, &membership); ok != nil {
		return nil, err
	}

	return &membership, err
}
