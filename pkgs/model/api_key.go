package model

type BaseAPIKey struct {
	Sys          *Sys            `json:"sys,omitempty"`
	Name         string          `json:"name,omitempty"`
	Description  string          `json:"description,omitempty"`
	AccessToken  string          `json:"accessToken,omitempty"`
	Environments []Environments  `json:"environments,omitempty"`
	Policies     []*APIKeyPolicy `json:"policies,omitempty"`
}

// APIKeyPolicy model
type APIKeyPolicy struct {
	Effect  string `json:"effect,omitempty"`
	Actions string `json:"actions,omitempty"`
}

// Environments model
type Environments struct {
	Sys Sys `json:"sys,omitempty"`
}

// PreviewAPIKey model
type PreviewAPIKey = BaseAPIKey

// APIKey model
type APIKey struct {
	BaseAPIKey
	PreviewAPIKey *BaseAPIKey `json:"preview_api_key,omitempty"`
}

// GetVersion returns entity version
func (apiKey *APIKey) GetVersion() int {
	version := 1
	if apiKey.Sys != nil {
		version = apiKey.Sys.Version
	}

	return version
}

func (apiKey *APIKey) IsNew() bool {
	return apiKey.Sys == nil || apiKey.Sys.ID == ""
}
