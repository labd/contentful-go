package contentful

import "fmt"

// EnvironmentsService service
type EnvironmentsService service

// Environment model
type Environment struct {
	Sys  *Sys   `json:"sys"`
	Name string `json:"name"`
}

// GetVersion returns entity version
func (e *Environment) GetVersion() int {
	version := 1
	if e.Sys != nil {
		version = e.Sys.Version
	}

	return version
}

// List returns an environments collection
func (service *EnvironmentsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments", spaceID)
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

// Get returns a single environment entity
func (service *EnvironmentsService) Get(spaceID, environmentID string) (*Environment, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s", spaceID, environmentID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var environment Environment
	if err := service.c.do(req, &environment); err != nil {
		return nil, err
	}

	return &environment, nil
}
