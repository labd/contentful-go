package model

import (
	"encoding/json"
	"io"
)

type Locale struct {
	Sys *EnvironmentSys `json:"sys,omitempty"`

	// Locale name
	Name string `json:"name,omitempty"`

	// Language code
	Code string `json:"code,omitempty"`

	// If no content is provided for the locale, the Delivery API will return content in a locale specified below:
	FallbackCode *string `json:"fallbackCode"`

	// Make the locale as default locale for your account
	Default bool `json:"default,omitempty"`

	// Entries with required fields can still be published if locale is empty.
	Optional bool `json:"optional,omitempty"`

	// Includes locale in the Delivery API response.
	CDA bool `json:"contentDeliveryApi"`

	// Displays locale to editors and enables it in Management API.
	CMA bool `json:"contentManagementApi"`
}

// GetVersion returns entity version
func (l *Locale) GetVersion() int {
	version := 1
	if l.Sys != nil {
		version = l.Sys.Version
	}

	return version
}

func (l *Locale) IsNew() bool {
	return l.Sys == nil || l.Sys.ID == ""
}

func (l *Locale) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(&l)
}
