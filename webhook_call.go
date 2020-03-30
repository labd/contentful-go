package contentful

import (
	"fmt"
	"net/http"
	"net/url"
)

// WebhookCallsService service
type WebhookCallsService service

// WebhookCall model
type WebhookCall struct {
	Sys        *Sys     `json:"sys"`
	Request    Request  `json:"request,omitempty"`
	Response   Response `json:"response,omitempty"`
	StatusCode int      `json:"statusCode"`
	Errors     []string `json:"errors"`
	EventType  string   `json:"eventType"`
	URL        string   `json:"url"`
	RequestAt  string   `json:"requestAt"`
	ResponseAt string   `json:"responseAt"`
}

// Request model
type Request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Response model
type Response struct {
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
}

// WebhookHealth model
type WebhookHealth struct {
	Sys   *Sys          `json:"sys"`
	Calls HealthDetails `json:"calls"`
}

// HealthDetails model
type HealthDetails struct {
	Total   int `json:"total"`
	Healthy int `json:"healthy"`
}

// List returns a webhook calls collection
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

// Get returns details of a single webhook call
func (service *WebhookCallsService) Get(spaceID, webhookID, callID string) (*WebhookCall, error) {
	path := fmt.Sprintf("/spaces/%s/webhooks/%s/calls/%s", spaceID, webhookID, callID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &WebhookCall{}, err
	}

	var webHook WebhookCall
	if ok := service.c.do(req, &webHook); ok != nil {
		return nil, err
	}

	return &webHook, err
}

// Health returns the health of a webhook
func (service *WebhookCallsService) Health(spaceID, webhookID string) (*WebhookHealth, error) {
	path := fmt.Sprintf("/spaces/%s/webhooks/%s/health", spaceID, webhookID)
	query := url.Values{}
	method := "GET"

	req, err := service.c.newRequest(method, path, query, nil)
	if err != nil {
		return &WebhookHealth{}, err
	}

	var health WebhookHealth
	if ok := service.c.do(req, &health); ok != nil {
		return nil, err
	}

	return &health, err
}
