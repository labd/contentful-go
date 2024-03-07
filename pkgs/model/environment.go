package model

import (
	"encoding/json"
	"io"
)

type StatusSys struct {
	SpaceSys
	Status *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"Status,omitempty"`

	AliasedEnvironment *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"aliasedEnvironment,omitempty"`

	Aliases []*struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"aliases,omitempty"`
}

type Environment struct {
	Sys  *StatusSys `json:"sys"`
	Name string     `json:"name"`
}

// GetVersion returns entity version
func (e *Environment) GetVersion() int {
	version := 1
	if e.Sys != nil {
		version = e.Sys.Version
	}

	return version
}

func (e *Environment) IsNew() bool {
	return e.Sys == nil || e.Sys.ID == ""
}

func (e *Environment) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(&e)
}
