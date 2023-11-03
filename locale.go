package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// LocalesService service
type LocalesService service

// Locale model
type Locale struct {
	Sys *Sys `json:"sys,omitempty"`

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
func (locale *Locale) GetVersion() int {
	version := 1
	if locale.Sys != nil {
		version = locale.Sys.Version
	}

	return version
}

// List returns a locales collection
func (service *LocalesService) List(spaceID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/locales", spaceID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}

// Get returns a single locale entity
func (service *LocalesService) Get(spaceID, localeID string) (*Locale, error) {
	path := fmt.Sprintf("/spaces/%s/locales/%s", spaceID, localeID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var locale Locale
	if err := service.c.do(req, &locale); err != nil {
		return nil, err
	}

	return &locale, nil
}

// GetWithEnv returns a single locale entity
func (service *LocalesService) GetWithEnv(env *Environment, localeID string) (*Locale, error) {
	path := fmt.Sprintf("/spaces/%s/environments/%s/locales/%s", env.Sys.Space.Sys.ID, env.Sys.ID, localeID)
	method := "GET"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return nil, err
	}

	var locale Locale
	if err := service.c.do(req, &locale); err != nil {
		return nil, err
	}

	return &locale, nil
}

// Delete the locale
func (service *LocalesService) Delete(spaceID string, locale *Locale) error {
	path := fmt.Sprintf("/spaces/%s/locales/%s", spaceID, locale.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(locale.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// DeleteWithEnv the locale
func (service *LocalesService) DeleteWithEnv(env *Environment, locale *Locale) error {
	path := fmt.Sprintf("/spaces/%s/environments/%s/locales/%s", env.Sys.Space.Sys.ID, env.Sys.ID, locale.Sys.ID)
	method := "DELETE"

	req, err := service.c.newRequest(method, path, nil, nil)
	if err != nil {
		return err
	}

	version := strconv.Itoa(locale.Sys.Version)
	req.Header.Set("X-Contentful-Version", version)

	return service.c.do(req, nil)
}

// UpsertWithEnv update or creates a locale for an environment
func (service *LocalesService) UpsertWithEnv(env *Environment, locale *Locale) error {
	var path string
	var method string

	path = fmt.Sprintf("/spaces/%s/environments/%s/locales", env.Sys.Space.Sys.ID, env.Sys.ID)
	method = "POST"
	if locale.Sys != nil && locale.Sys.CreatedAt != "" {
		path = fmt.Sprintf("%s/%s", path, locale.Sys.ID)
		method = "PUT"
	}

	return service.doUpsert(path, method, locale)
}

// Upsert updates or creates a new locale entity
func (service *LocalesService) Upsert(spaceID string, locale *Locale) error {
	var path string
	var method string

	path = fmt.Sprintf("/spaces/%s/locales", spaceID)
	method = "POST"

	if locale.Sys != nil && locale.Sys.CreatedAt != "" {
		path = fmt.Sprintf("%s/%s", path, locale.Sys.ID)
		method = "PUT"
	}

	return service.doUpsert(path, method, locale)
}

func (service *LocalesService) doUpsert(path string, method string, locale *Locale) error {
	bytesArray, err := json.Marshal(locale)
	if err != nil {
		return err
	}

	req, err := service.c.newRequest(method, path, nil, bytes.NewReader(bytesArray))
	if err != nil {
		return err
	}

	req.Header.Set("X-Contentful-Version", strconv.Itoa(locale.GetVersion()))

	return service.c.do(req, locale)
}
