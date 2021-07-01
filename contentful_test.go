package contentful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	server         *httptest.Server
	cma            *Client
	urc            *Client
	c              *Client
	CMAToken       = "b4c0n73n7fu1"
	CDAToken       = "cda-token"
	CPAToken       = "cpa-token"
	spaceID        = "id1"
	organizationID = "org-id"
	environmentID  = "env-id"
	env            = &Environment{
		Sys: &Sys{
			ID:       environmentID,
			LinkType: "env-link-type",
			Type:     "env-type",
			Space: &Space{
				Name:          "space-name",
				DefaultLocale: "en",
				Sys: &Sys{
					ID:       spaceID,
					LinkType: "space-link-type",
					Type:     "space-type",
				},
			},
		},
		Name: "env-name",
	}
)

func readTestData(fileName string) string {
	path := "testdata/" + fileName
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}

func checkHeaders(req *http.Request, assert *assert.Assertions) {
	assert.Equal("Bearer "+CMAToken, req.Header.Get("Authorization"))
	assert.Equal("application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
}

func spaceFromTestData(fileName string) (*Space, error) {
	content := readTestData(fileName)

	var space Space
	err := json.NewDecoder(strings.NewReader(content)).Decode(&space)
	if err != nil {
		return nil, err
	}

	return &space, nil
}

func webhookFromTestData(fileName string) (*Webhook, error) {
	content := readTestData(fileName)

	var webhook Webhook
	err := json.NewDecoder(strings.NewReader(content)).Decode(&webhook)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

func contentTypeFromTestData(fileName string) (*ContentType, error) {
	content := readTestData(fileName)

	var ct ContentType
	err := json.NewDecoder(strings.NewReader(content)).Decode(&ct)
	if err != nil {
		return nil, err
	}

	return &ct, nil
}

func localeFromTestData(fileName string) (*Locale, error) {
	content := readTestData(fileName)

	var locale Locale
	err := json.NewDecoder(strings.NewReader(content)).Decode(&locale)
	if err != nil {
		return nil, err
	}

	return &locale, nil
}

func environmentFromTestData(fileName string) (*Environment, error) {
	content := readTestData(fileName)

	var environment Environment
	err := json.NewDecoder(strings.NewReader(content)).Decode(&environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func environmentAliasFromTestData(fileName string) (*EnvironmentAlias, error) {
	content := readTestData(fileName)

	var environmentAlias EnvironmentAlias
	err := json.NewDecoder(strings.NewReader(content)).Decode(&environmentAlias)
	if err != nil {
		return nil, err
	}

	return &environmentAlias, nil
}

func entryFromTestData(fileName string) (*Entry, error) {
	content := readTestData(fileName)

	var entry Entry
	err := json.NewDecoder(strings.NewReader(content)).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func roleFromTestData(fileName string) (*Role, error) {
	content := readTestData(fileName)

	var role Role
	err := json.NewDecoder(strings.NewReader(content)).Decode(&role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func membershipFromTestData(fileName string) (*Membership, error) {
	content := readTestData(fileName)

	var membership Membership
	err := json.NewDecoder(strings.NewReader(content)).Decode(&membership)
	if err != nil {
		return nil, err
	}

	return &membership, nil
}

func assetFromTestData(fileName string) (*Asset, error) {
	content := readTestData(fileName)

	var asset Asset
	err := json.NewDecoder(strings.NewReader(content)).Decode(&asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func apiKeyFromTestData(fileName string) (*APIKey, error) {
	content := readTestData(fileName)

	var apiKey APIKey
	err := json.NewDecoder(strings.NewReader(content)).Decode(&apiKey)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func accessTokenFromTestFile(fileName string) (*AccessToken, error) {
	content := readTestData(fileName)

	var accessToken AccessToken
	err := json.NewDecoder(strings.NewReader(content)).Decode(&accessToken)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func entryTaskFromTestFile(fileName string) (*EntryTask, error) {
	content := readTestData(fileName)

	var entryTask EntryTask
	err := json.NewDecoder(strings.NewReader(content)).Decode(&entryTask)
	if err != nil {
		return nil, err
	}

	return &entryTask, nil
}

func scheduledActionFromTestFile(fileName string) (*ScheduledAction, error) {
	content := readTestData(fileName)

	var scheduledAction ScheduledAction
	err := json.NewDecoder(strings.NewReader(content)).Decode(&scheduledAction)
	if err != nil {
		return nil, err
	}

	return &scheduledAction, nil
}

func editorInterfaceFromTestFile(fileName string) (*EditorInterface, error) {
	content := readTestData(fileName)

	var editorInterface EditorInterface
	err := json.NewDecoder(strings.NewReader(content)).Decode(&editorInterface)
	if err != nil {
		return nil, err
	}

	return &editorInterface, nil
}

func extensionFromTestFile(fileName string) (*Extension, error) {
	content := readTestData(fileName)

	var extension Extension
	err := json.NewDecoder(strings.NewReader(content)).Decode(&extension)
	if err != nil {
		return nil, err
	}

	return &extension, nil
}

func appDefinitionFromTestFile(fileName string) (*AppDefinition, error) {
	content := readTestData(fileName)

	var appDefinition AppDefinition
	err := json.NewDecoder(strings.NewReader(content)).Decode(&appDefinition)
	if err != nil {
		return nil, err
	}

	return &appDefinition, nil
}

func appInstallationFromTestFile(fileName string) (*AppInstallation, error) {
	content := readTestData(fileName)

	var appInstallation AppInstallation
	err := json.NewDecoder(strings.NewReader(content)).Decode(&appInstallation)
	if err != nil {
		return nil, err
	}

	return &appInstallation, nil
}

func resourceFromTestFile(fileName string) (*Resource, error) {
	content := readTestData(fileName)

	var resource Resource
	err := json.NewDecoder(strings.NewReader(content)).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func setup() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture := strings.Replace(r.URL.Path, "/", "-", -1)
		fixture = strings.TrimLeft(fixture, "-")
		var path string

		if e := r.URL.Query().Get("error"); e != "" {
			path = "testdata/error_" + e + ".json"
		} else {
			if r.Method == "GET" {
				path = "testdata/" + fixture + ".json"
			}

			if r.Method == "POST" {
				path = "testdata/" + fixture + "_new.json"
			}

			if r.Method == "PUT" {
				path = "testdata/" + fixture + "_updated.json"
			}
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			_, _ = fmt.Fprintln(w, err)
			return
		}

		_, _ = fmt.Fprintln(w, string(file))
		return
	})

	server = httptest.NewServer(handler)

	c = NewCMA(CMAToken)
	c.BaseURL = server.URL
}

func teardown() {
	server.Close()
	c = nil
}

func TestContentfulNewCMA(t *testing.T) {
	assertions := assert.New(t)

	cma := NewCMA(CMAToken)
	assertions.IsType(Client{}, *cma)
	assertions.Equal("https://api.contentful.com", cma.BaseURL)
	assertions.Equal("CMA", cma.api)
	assertions.Equal(CMAToken, cma.token)
	assertions.Equal(fmt.Sprintf("Bearer %s", CMAToken), cma.Headers["Authorization"])
	assertions.Equal("application/vnd.contentful.management.v1+json", cma.Headers["Content-Type"])
	assertions.Equal(fmt.Sprintf("sdk contentful.go/%s", Version), cma.Headers["X-Contentful-User-Agent"])
}

func TestContentfulNewCDA(t *testing.T) {
	assertions := assert.New(t)

	cda := NewCDA(CDAToken)
	assertions.IsType(Client{}, *cda)
	assertions.Equal("https://cdn.contentful.com", cda.BaseURL)
	assertions.Equal("CDA", cda.api)
	assertions.Equal(CDAToken, cda.token)
	assertions.Equal(fmt.Sprintf("Bearer %s", CDAToken), cda.Headers["Authorization"])
	assertions.Equal("application/vnd.contentful.delivery.v1+json", cda.Headers["Content-Type"])
	assertions.Equal(fmt.Sprintf("contentful-go/%s", Version), cda.Headers["X-Contentful-User-Agent"])
}

func TestContentfulNewCPA(t *testing.T) {
	assertions := assert.New(t)

	cpa := NewCPA(CPAToken)
	assertions.IsType(Client{}, *cpa)
	assertions.Equal("https://preview.contentful.com", cpa.BaseURL)
	assertions.Equal("CPA", cpa.api)
	assertions.Equal(CPAToken, cpa.token)
}

func TestNewResourceClient(t *testing.T) {
	assertions := assert.New(t)

	urc := NewResourceClient(CMAToken)
	assertions.IsType(Client{}, *urc)
	assertions.Equal("https://upload.contentful.com", urc.BaseURL)
	assertions.Equal("URC", urc.api)
	assertions.Equal(CMAToken, urc.token)
	assertions.Equal(fmt.Sprintf("Bearer %s", CMAToken), urc.Headers["Authorization"])
}

func TestContentfulSetOrganization(t *testing.T) {
	assertions := assert.New(t)

	cma := NewCMA(CMAToken)
	cma.SetOrganization(organizationID)
	assertions.Equal(organizationID, cma.Headers["X-Contentful-Organization"])
}

func TestContentfulSetClient(t *testing.T) {
	assertions := assert.New(t)

	newClient := &http.Client{}
	cma := NewCMA(CMAToken)
	cma.SetHTTPClient(newClient)
	assertions.Equal(newClient, cma.client)
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	assertions := assert.New(t)

	method := "GET"
	path := "/some/path"
	query := url.Values{}
	query.Add("foo", "bar")
	query.Add("faz", "zoo")

	expectedURL, _ := url.Parse(c.BaseURL)
	expectedURL.Path = path
	expectedURL.RawQuery = query.Encode()

	req, err := c.newRequest(method, path, query, nil)
	assertions.Nil(err)
	assertions.Equal(req.Header.Get("Authorization"), "Bearer "+CMAToken)
	assertions.Equal(req.Header.Get("Content-Type"), "application/vnd.contentful.management.v1+json")
	assertions.Equal(req.Method, method)
	assertions.Equal(req.URL.String(), expectedURL.String())

	method = "POST"
	type RequestBody struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	bodyData := RequestBody{
		Name: "test",
		Age:  10,
	}
	body, _ := json.Marshal(bodyData)
	req, err = c.newRequest(method, path, query, bytes.NewReader(body))
	assertions.Nil(err)
	assertions.Equal(req.Header.Get("Authorization"), "Bearer "+CMAToken)
	assertions.Equal(req.Header.Get("Content-Type"), "application/vnd.contentful.management.v1+json")
	assertions.Equal(req.Method, method)
	assertions.Equal(req.URL.String(), expectedURL.String())
	defer req.Body.Close()
	var requestBody RequestBody
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	assertions.Nil(err)
	assertions.Equal(requestBody, bodyData)
}

func TestHandleError(t *testing.T) {
	setup()
	defer teardown()

	assertions := assert.New(t)

	method := "GET"
	path := "/some/path"
	requestID := "request-id"
	query := url.Values{}
	errResponse := ErrorResponse{
		Sys: &Sys{
			ID:   "AccessTokenInvalid",
			Type: "Error",
		},
		Message:   "Access token is invalid",
		RequestID: requestID,
	}

	marshaled, _ := json.Marshal(errResponse)
	errResponseReader := bytes.NewReader(marshaled)
	errResponseReadCloser := ioutil.NopCloser(errResponseReader)

	req, _ := c.newRequest(method, path, query, nil)
	responseHeaders := http.Header{}
	responseHeaders.Add("X-Contentful-Request-Id", requestID)
	res := &http.Response{
		Header:     responseHeaders,
		StatusCode: http.StatusUnauthorized,
		Body:       errResponseReadCloser,
		Request:    req,
	}

	err := c.handleError(req, res)
	assertions.IsType(AccessTokenInvalidError{}, err)
	assertions.Equal(req, err.(AccessTokenInvalidError).APIError.req)
	assertions.Equal(res, err.(AccessTokenInvalidError).APIError.res)
	assertions.Equal(&errResponse, err.(AccessTokenInvalidError).APIError.err)
}

func TestBackoffForPerSecondLimiting(t *testing.T) {
	var err error
	assertions := assert.New(t)
	rateLimited := true
	waitSeconds := 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rateLimited == true {
			w.Header().Set("X-Contentful-Request-Id", "request-id")
			w.Header().Set("Content-Type", "application/vnd.contentful.management.v1+json")
			w.Header().Set("X-Contentful-Ratelimit-Hour-Limit", "36000")
			w.Header().Set("X-Contentful-Ratelimit-Hour-Remaining", "35883")
			w.Header().Set("X-Contentful-Ratelimit-Reset", strconv.Itoa(waitSeconds))
			w.Header().Set("X-Contentful-Ratelimit-Second-Limit", "10")
			w.Header().Set("X-Contentful-Ratelimit-Second-Remaining", "0")
			w.WriteHeader(429)

			_, _ = w.Write([]byte(readTestData("error_ratelimit.json")))
		} else {
			_, _ = w.Write([]byte(readTestData("space-1.json")))
		}
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	go func() {
		time.Sleep(time.Second * time.Duration(waitSeconds))
		rateLimited = false
	}()

	space, err := cma.Spaces.Get("id1")
	assertions.Nil(err)
	assertions.Equal(space.Name, "Contentful Example API")
	assertions.Equal(space.Sys.ID, "id1")
}
