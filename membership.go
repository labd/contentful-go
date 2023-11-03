package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// MembershipsService service
type MembershipsService service

// Membership model
type Membership struct {
	Sys   *Sys    `json:"sys"`
	Admin bool    `json:"admin"`
	Roles []Roles `json:"roles"`
	User  Member  `json:"user,omitempty"`
	Email string  `json:"email,omitempty"`
}

// Roles model
type Roles struct {
	Sys *Sys `json:"sys"`
}

// Member model
type Member struct {
	Sys *Sys `json:"sys"`
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

// Upsert updates or creates a new membership
func (service *MembershipsService) Upsert(spaceID string, m *Membership) error {
	bytesArray, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var path string
	var method string

	if m.Sys != nil && m.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s/space_memberships/%s", spaceID, m.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/space_memberships", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(m.GetVersion()))

	return service.c.do(req, m)
}

// Delete the role
func (service *MembershipsService) Delete(spaceID string, membershipID string) error {
	path := fmt.Sprintf("/spaces/%s/space_memberships/%s", spaceID, membershipID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	return service.c.do(req, nil)
}
