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

// ToSpace cast Items to Space model
func (col *Collection) ToSpace() []*Space {
	var spaces []*Space

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&spaces)

	return spaces
}

// ToWebhook cast Items to Webhook model
func (col *Collection) ToWebhook() []*Webhook {
	var webhooks []*Webhook

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&webhooks)

	return webhooks
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

// ToAccessToken cast Items to AccessToken model
func (col *Collection) ToAccessToken() []*AccessToken {
	var accessTokens []*AccessToken

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&accessTokens)

	return accessTokens
}

// ToEntryTask cast Items to EntryTask model
func (col *Collection) ToEntryTask() []*EntryTask {
	var entryTasks []*EntryTask

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&entryTasks)

	return entryTasks
}

// ToScheduledAction cast Items to ScheduledAction model
func (col *Collection) ToScheduledAction() []*ScheduledAction {
	var scheduledActions []*ScheduledAction

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&scheduledActions)

	return scheduledActions
}

// ToEditorInterface cast Items to EditorInterface model
func (col *Collection) ToEditorInterface() []*EditorInterface {
	var editorInterface []*EditorInterface

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&editorInterface)

	return editorInterface
}

// ToExtension cast Items to Extension model
func (col *Collection) ToExtension() []*Extension {
	var extension []*Extension

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&extension)

	return extension
}

// ToWebhookCall cast Items to WebhookCall model
func (col *Collection) ToWebhookCall() []*WebhookCall {
	var webhookCall []*WebhookCall

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&webhookCall)

	return webhookCall
}

// ToUsage cast Items to Usage model
func (col *Collection) ToUsage() []*Usage {
	var usage []*Usage

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&usage)

	return usage
}

// ToMembership cast Items to Membership model
func (col *Collection) ToMembership() []*Membership {
	var membership []*Membership

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&membership)

	return membership
}

// ToRole cast Items to Role model
func (col *Collection) ToRole() []*Role {
	var role []*Role

	byteArray, _ := json.Marshal(col.Items)
	_ = json.NewDecoder(bytes.NewReader(byteArray)).Decode(&role)

	return role
}
