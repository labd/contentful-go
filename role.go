package contentful

import (
	"fmt"
)

// RolesService service
type RolesService service

// Role model
type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Policies    []struct {
		Effect     string   `json:"effect"`
		Actions    []string `json:"actions"`
		Constraint struct {
			And []struct {
				Equals []interface{} `json:"equals"`
			} `json:"and"`
		} `json:"constraint"`
	} `json:"policies"`
	Permissions struct {
		ContentModel       []string      `json:"ContentModel"`
		Settings           []interface{} `json:"Settings"`
		ContentDelivery    []interface{} `json:"ContentDelivery"`
		Environments       []interface{} `json:"Environments"`
		EnvironmentAliases []interface{} `json:"EnvironmentAliases"`
	} `json:"permissions"`
	Sys *Sys `json:"sys"`
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
