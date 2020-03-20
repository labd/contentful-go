package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// RolesService service
type RolesService service

// Role model
type Role struct {
	Sys         *Sys        `json:"sys"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Policies    []Policies  `json:"policies"`
	Permissions Permissions `json:"permissions"`
}

type Policies struct {
	Effect     string     `json:"effect"`
	Actions    []string   `json:"actions"`
	Constraint Constraint `json:"constraint"`
}

type Permissions struct {
	ContentModel       []string `json:"ContentModel"`
	Settings           string   `json:"Settings"`
	ContentDelivery    string   `json:"ContentDelivery"`
	Environments       string   `json:"Environments"`
	EnvironmentAliases string   `json:"EnvironmentAliases"`
}

type Constraint struct {
	And []ConstraintDetail `json:"and"`
}

type ConstraintDetail struct {
	Equals DetailItem `json:"equals"`
}

type DetailItem struct {
	Doc      map[string]interface{}
	ItemType string
}

// GetVersion returns entity version
func (r *Role) GetVersion() int {
	version := 1
	if r.Sys != nil {
		version = r.Sys.Version
	}

	return version
}

// List returns an environments collection
func (service *RolesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/roles", spaceID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single role
func (service *RolesService) Get(spaceID, roleID string) (*Role, error) {
	path := fmt.Sprintf("/spaces/%s/roles/%s", spaceID, roleID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &Role{}, err
	}

	var role Role
	if ok := service.c.do(req, &role); ok != nil {
		return nil, err
	}

	return &role, err
}

// Upsert updates or creates a new role
func (service *RolesService) Upsert(spaceID string, r *Role) error {
	bytesArray, err := json.Marshal(r)
	fmt.Println(r)
	if err != nil {
		return err
	}

	var path string
	var method string

	if r.Sys != nil && r.Sys.ID != "" {
		path = fmt.Sprintf("/spaces/%s/roles/%s", spaceID, r.Sys.ID)
		method = "PUT"
	} else {
		path = fmt.Sprintf("/spaces/%s/roles", spaceID)
		method = "POST"
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(r.GetVersion()))

	return service.c.do(req, r)
}
