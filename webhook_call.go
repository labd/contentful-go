package contentful

import (
	"fmt"
	"net/http"
)

// WebhookCallsService service
type WebhookCallsService service

// WebhookCall model
type WebhookCall struct {
	Sys        *Sys     `json:"sys"`
	StatusCode int      `json:"statusCode"`
	Errors     []string `json:"errors"`
	EventType  string   `json:"eventType"`
	URL        string   `json:"url"`
	RequestAt  string   `json:"requestAt"`
	ResponseAt string   `json:"responseAt"`
}

// GetVersion returns entity version
func (webhookCall *WebhookCall) GetVersion() int {
	version := 1
	if webhookCall.Sys != nil {
		version = webhookCall.Sys.Version
	}

	return version
}

// List returns entries collection
func (service *WebhookCallsService) List(spaceID, webhookID string) *Collection {
	path := fmt.Sprintf("/spaces/%s/webhooks/%s/calls", spaceID, webhookID)

	req, err := service.c.newRequest(http.MethodGet, path, nil, nil)
	if err != nil {
		return &Collection{}
	}

	col := NewCollection(&CollectionOptions{})
	col.c = service.c
	col.req = req

	return col
}
