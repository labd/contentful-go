package contentful

import (
	"fmt"
)

// ExtensionsService service
type ExtensionsService service

// Extension model
type Extension struct {
	Sys       *Sys             `json:"sys"`
	Extension ExtensionDetails `json:"extension"`
}

// ExtensionDetails model
type ExtensionDetails struct {
	SRC        string                 `json:"src"`
	Name       string                 `json:"name"`
	FieldTypes map[string]interface{} `json:"fieldTypes"`
	Sidebar    bool                   `json:"sidebar"`
}

// List returns an extensions collection
func (service *ExtensionsService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/extensions", spaceID, service.c.Environment)

	req, err := service.c.newRequest("GET", path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}
