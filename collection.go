package contentful

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// CollectionOptions holds init options
type CollectionOptions struct {
	Limit uint16
}

// Collection model
type Collection struct {
	Query
	c        *Client
	req      *http.Request
	page     uint16
	Sys      *Sys          `json:"sys"`
	Total    int           `json:"total"`
	Skip     int           `json:"skip"`
	Limit    int           `json:"limit"`
	Items    []interface{} `json:"items"`
	Includes interface{}   `json:"includes"`
}

// NewCollection initializes a new collection
func NewCollection(options *CollectionOptions) *Collection {
	query := NewQuery()
	query.Order("sys.createdAt", true)

	if options.Limit > 0 {
		query.Limit(options.Limit)
	}

	return &Collection{
		Query: *query,
		page:  1,
	}
}

// Next makes the col.req
func (col *Collection) Next() (*Collection, error) {
	// setup query params
	skip := uint16(col.Limit) * (col.page - 1)
	col.Query.Skip(skip)

	// override request query
	col.req.URL.RawQuery = col.Query.String()

	// makes api call
	err := col.c.do(col.req, col)
	if err != nil {
		return nil, err
	}

	col.page++

	return col, nil
}

// ToContentType cast Items to ContentType model
func (col *Collection) ToContentType() []*ContentType {
	var contentTypes []*ContentType

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&contentTypes)

	return contentTypes
}

// ToEnvironment cast Items to Environment model
func (col *Collection) ToEnvironment() []*Environment {
	var environments []*Environment

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&environments)

	return environments
}

// ToEnvironmentAlias cast Items to EnvironmentAlias model
func (col *Collection) ToEnvironmentAlias() []*EnvironmentAlias {
	var environmentAlias []*EnvironmentAlias

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&environmentAlias)

	return environmentAlias
}

// ToSpace cast Items to Space model
func (col *Collection) ToSpace() []*Space {
	var spaces []*Space

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&spaces)

	return spaces
}

// ToEntry cast Items to Entry model
func (col *Collection) ToEntry() []*Entry {
	var entries []*Entry

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&entries)

	return entries
}

// ToLocale cast Items to Locale model
func (col *Collection) ToLocale() []*Locale {
	var locales []*Locale

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&locales)

	return locales
}

// ToAsset cast Items to Asset model
func (col *Collection) ToAsset() []*Asset {
	var assets []*Asset

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&assets)

	return assets
}

// ToAPIKey cast Items to APIKey model
func (col *Collection) ToAPIKey() []*APIKey {
	var apiKeys []*APIKey

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&apiKeys)

	return apiKeys
}

// ToWebhook cast Items to Webhook model
func (col *Collection) ToWebhook() []*Webhook {
	var webhooks []*Webhook

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&webhooks)

	return webhooks
}

// ToUser cast Items to User model
func (col *Collection) ToUser() []*User {
	var user []*User

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&user)

	return user
}

// ToOrganization cast Items to Organization model
func (col *Collection) ToOrganization() []*Organization {
	var organization []*Organization

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&organization)

	return organization
}

// ToEntrySnapshot cast Items to Snapshot model for entries
func (col *Collection) ToEntrySnapshot() []*EntrySnapshot {
	var snapshot []*EntrySnapshot

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&snapshot)

	return snapshot
}

// ToContentTypeSnapshot cast Items to Snapshot model for content types
func (col *Collection) ToContentTypeSnapshot() []*ContentTypeSnapshot {
	var snapshot []*ContentTypeSnapshot

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&snapshot)

	return snapshot
}

// ToAccessToken cast Items to AccessToken model for content types
func (col *Collection) ToAccessToken() []*AccessToken {
	var accessTokens []*AccessToken

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&accessTokens)

	return accessTokens
}
