package contentful

import "fmt"

// EnvironmentAliasesService service
type EnvironmentAliasesService service

// EnvironmentAlias model
type EnvironmentAlias struct {
	Sys              *Sys         `json:"sys"`
	EnvironmentAlias *AliasDetail `json:"environment"`
}

// AliasDetail model
type AliasDetail struct {
	Sys *Sys `json:"sys"`
}

// GetVersion returns entity version
func (environmentAlias *EnvironmentAlias) GetVersion() int {
	version := 1
	if environmentAlias.Sys != nil {
		version = environmentAlias.Sys.Version
	}

	return version
}

// List returns an environment aliases collection
func (service *EnvironmentAliasesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environment_aliases", spaceID)
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

// Get returns a single environment alias entity
func (service *EnvironmentAliasesService) Get(spaceID, environmentAliasID string) (*EnvironmentAlias, error) {
	path := fmt.Sprintf("/spaces/%s/environment_aliases/%s", spaceID, environmentAliasID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var environmentAlias EnvironmentAlias
	if err := service.c.do(req, &environmentAlias); err != nil {
		return nil, err
	}

	return &environmentAlias, nil
}
