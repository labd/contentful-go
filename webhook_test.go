package contentful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhooksService_List(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/webhook_definitions")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("webhook.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Webhooks.List(spaceID).Next()
	assertions.Nil(err)
	webhook := collection.ToWebhook()
	assertions.Equal(1, len(webhook))
	assertions.Equal("webhook-name", webhook[0].Name)
}

func TestWebhooksService_Get(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "GET")
		assertions.Equal(r.URL.Path, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI")

		checkHeaders(r, assertions)

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("webhook_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	webhook, err := cma.Webhooks.Get(spaceID, "7fstd9fZ9T2p3kwD49FxhI")
	assertions.Nil(err)
	assertions.Equal("webhook-name", webhook.Name)
}

func TestWebhooksService_Upsert_Create(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "POST")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("webhook-name", payload["name"])
		assertions.Equal("https://www.example.com/test", payload["url"])
		assertions.Equal("username", payload["httpBasicUsername"])
		assertions.Equal("password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assertions.Equal(2, len(topics))
		assertions.Equal("Entry.create", topics[0].(string))
		assertions.Equal("ContentType.create", topics[1].(string))

		headers := payload["headers"].([]interface{})
		assertions.Equal(2, len(headers))
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assertions.Equal("header1", header1["key"].(string))
		assertions.Equal("header1-value", header1["value"].(string))

		assertions.Equal("header2", header2["key"].(string))
		assertions.Equal("header2-value", header2["value"].(string))

		w.WriteHeader(201)
		_, _ = fmt.Fprintln(w, readTestData("webhook_1.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	webhook := &Webhook{
		Name: "webhook-name",
		URL:  "https://www.example.com/test",
		Topics: []string{
			"Entry.create",
			"ContentType.create",
		},
		HTTPBasicUsername: "username",
		HTTPBasicPassword: "password",
		Headers: []*WebhookHeader{
			{
				Key:   "header1",
				Value: "header1-value",
			},
			{
				Key:   "header2",
				Value: "header2-value",
			},
		},
	}

	err = cma.Webhooks.Upsert(spaceID, webhook)
	assertions.Nil(err)
	assertions.Equal("7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assertions.Equal("webhook-name", webhook.Name)
	assertions.Equal("username", webhook.HTTPBasicUsername)
}

func TestWebhooksService_Upsert_Update(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "PUT")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI")
		checkHeaders(r, assertions)

		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assertions.Nil(err)
		assertions.Equal("updated-webhook-name", payload["name"])
		assertions.Equal("https://www.example.com/test-updated", payload["url"])
		assertions.Equal("updated-username", payload["httpBasicUsername"])
		assertions.Equal("updated-password", payload["httpBasicPassword"])

		topics := payload["topics"].([]interface{})
		assertions.Equal(3, len(topics))
		assertions.Equal("Entry.create", topics[0].(string))
		assertions.Equal("ContentType.create", topics[1].(string))
		assertions.Equal("Asset.create", topics[2].(string))

		headers := payload["headers"].([]interface{})
		assertions.Equal(2, len(headers))
		header1 := headers[0].(map[string]interface{})
		header2 := headers[1].(map[string]interface{})

		assertions.Equal("header1", header1["key"].(string))
		assertions.Equal("updated-header1-value", header1["value"].(string))

		assertions.Equal("header2", header2["key"].(string))
		assertions.Equal("updated-header2-value", header2["value"].(string))

		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, readTestData("webhook_updated.json"))
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData("webhook_1.json")
	assertions.Nil(err)

	webhook.Name = "updated-webhook-name"
	webhook.URL = "https://www.example.com/test-updated"
	webhook.Topics = []string{
		"Entry.create",
		"ContentType.create",
		"Asset.create",
	}
	webhook.HTTPBasicUsername = "updated-username"
	webhook.HTTPBasicPassword = "updated-password"
	webhook.Headers = []*WebhookHeader{
		{
			Key:   "header1",
			Value: "updated-header1-value",
		},
		{
			Key:   "header2",
			Value: "updated-header2-value",
		},
	}

	err = cma.Webhooks.Upsert(spaceID, webhook)
	assertions.Nil(err)
	assertions.Equal("7fstd9fZ9T2p3kwD49FxhI", webhook.Sys.ID)
	assertions.Equal(1, webhook.Sys.Version)
	assertions.Equal("updated-webhook-name", webhook.Name)
	assertions.Equal("updated-username", webhook.HTTPBasicUsername)
}

func TestWebhookDelete(t *testing.T) {
	var err error
	assertions := assert.New(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertions.Equal(r.Method, "DELETE")
		assertions.Equal(r.RequestURI, "/spaces/"+spaceID+"/webhook_definitions/7fstd9fZ9T2p3kwD49FxhI")
		checkHeaders(r, assertions)

		w.WriteHeader(200)
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	// test webhook
	webhook, err := webhookFromTestData("webhook_1.json")
	assertions.Nil(err)

	err = cma.Webhooks.Delete(spaceID, webhook)
	assertions.Nil(err)
}
