package model

type AppDefinitionSys struct {
	CreatedSys
	Shared       bool `json:"shared,omitempty"`
	Organization *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"organization,omitempty"`
}

// AppDefinition model
type AppDefinition struct {
	Sys       *CreatedSys `json:"sys"`
	Name      string      `json:"name"`
	SRC       *string     `json:"src,omitempty"`
	Bundle    *Bundle     `json:"bundle,omitempty"`
	Locations []Locations `json:"locations"`
}

type Bundle struct {
	Sys *BaseSys `json:"sys"`
}

// Locations model
type Locations struct {
	Location       string          `json:"location"`
	FieldTypes     []FieldType     `json:"fieldTypes,omitempty"`
	NavigationItem *NavigationItem `json:"navigationItem,omitempty"`
}

type FieldType struct {
	Type     string  `json:"type"`
	LinkType *string `json:"linkType,omitempty"`
	Items    *Items  `json:"items,omitempty"`
}

type NavigationItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Items struct {
	Type     string  `json:"type"`
	LinkType *string `json:"linkType,omitempty"`
}

func (appDefinition *AppDefinition) IsNew() bool {
	return appDefinition.Sys == nil || appDefinition.Sys.ID == ""
}
