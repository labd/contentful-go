package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// EditorInterfacesService service
type EditorInterfacesService service

// EditorInterface model
type EditorInterface struct {
	Sys      *Sys       `json:"sys"`
	Controls []Controls `json:"controls"`
	SideBar  []Sidebar  `json:"sidebar,omitempty"`
}

// Controls model
type Controls struct {
	FieldID         string    `json:"fieldId"`
	WidgetNameSpace *string   `json:"widgetNamespace,omitempty"`
	WidgetID        *string   `json:"widgetId,omitempty"`
	Settings        *Settings `json:"settings,omitempty"`
}

type Settings struct {
	HelpText        *string `json:"helpText,omitempty"`
	TrueLabel       *string `json:"trueLabel,omitempty"`
	FalseLabel      *string `json:"falseLabel,omitempty"`
	Stars           *int64  `json:"stars,omitempty"`
	Format          *string `json:"format,omitempty"`
	AMPM            *string `json:"ampm,omitempty"`
	BulkEditing     *bool   `json:"bulkEditing,omitempty"`
	TrackingFieldId *string `json:"trackingFieldId,omitempty"`
}

// Sidebar model
type Sidebar struct {
	WidgetNameSpace string            `json:"widgetNamespace"`
	WidgetID        string            `json:"widgetId"`
	Settings        map[string]string `json:"settings,omitempty"`
	Disabled        bool              `json:"disabled"`
}

// List returns an EditorInterface collection
func (service *EditorInterfacesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/environments/%s/editor_interface", spaceID, service.c.Environment)

	req, err := service.c.newRequest("GET", path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single EditorInterface
func (service *EditorInterfacesService) Get(spaceID, contentTypeID string) (*EditorInterface, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/editor_interface", spaceID, service.c.Environment, contentTypeID)
	return service.doGet(path)
}

// GetWithEnv returns a single EditorInterface
func (service *EditorInterfacesService) GetWithEnv(env *Environment, contentTypeID string) (*EditorInterface, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/editor_interface", env.Sys.Space.Sys.ID, env.Sys.ID, contentTypeID)
	return service.doGet(path)
}

func (service *EditorInterfacesService) doGet(path string) (*EditorInterface, error) {
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &EditorInterface{}, err
	}

	var editorInterface EditorInterface
	if err = service.c.do(req, &editorInterface); err != nil {
		return nil, err
	}

	return &editorInterface, err
}

// Update updates an editor interface
func (service *EditorInterfacesService) Update(spaceID, contentTypeID string, e *EditorInterface) error {
	bytesArray, err := json.Marshal(e)
	if err != nil {
		return err
	}

	var path string
	var method string

	path = fmt.Sprintf("/spaces/%s/environments/%s/content_types/%s/editor_interface", spaceID, service.c.Environment, contentTypeID)
	method = "PUT"

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(e.Sys.Version))

	return service.c.do(req, e)
}
