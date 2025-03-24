package model

type AppInstallationSys struct {
	EnvironmentSys
	AppDefinition *struct {
		Sys BaseSys `json:"sys,omitempty"`
	} `json:"appDefinition,omitempty"`
}

type AppInstallation struct {
	Sys        *AppInstallationSys `json:"sys,omitempty"`
	Parameters map[string]any      `json:"parameters"`
	Terms      []string            `json:"-"`
}
