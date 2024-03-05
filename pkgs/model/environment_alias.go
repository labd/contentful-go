package model

import (
	"encoding/json"
	"io"
)

type EnvironmentAlias struct {
	Sys   *CreatedSys `json:"sys"`
	Alias *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"environment"`
}

// GetVersion returns entity version
func (e *EnvironmentAlias) GetVersion() int {
	version := 1
	if e.Sys != nil {
		version = e.Sys.Version
	}

	return version
}

func (e *EnvironmentAlias) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(&e)
}
