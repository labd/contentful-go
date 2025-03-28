package contentful

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/labd/contentful-go/internal/cda"
	"github.com/labd/contentful-go/internal/cma"
	"github.com/labd/contentful-go/pkgs/client"
	"github.com/labd/contentful-go/pkgs/common"
	cda_service "github.com/labd/contentful-go/service/cda"
	cma_service "github.com/labd/contentful-go/service/cma"

	"moul.io/http2curl"
)

// Client model
type Client struct {
	client        *http.Client
	api           string
	token         string
	Debug         bool
	QueryParams   map[string]string
	Headers       map[string]string
	BaseURL       string
	Environment   string
	commonService service

	Spaces           *SpacesService
	Users            *UsersService
	Organizations    *OrganizationsService
	Roles            *RolesService
	Memberships      *MembershipsService
	Snapshots        *SnapshotsService
	AccessTokens     *AccessTokensService
	EntryTasks       *EntryTasksService
	ScheduledActions *ScheduledActionsService
	Webhooks         *WebhooksService
	WebhookCalls     *WebhookCallsService
	EditorInterfaces *EditorInterfacesService
	Extensions       *ExtensionsService
	Usages           *UsagesService
	Resources        *ResourcesService
	AppUpload        *AppUploadService
	AppBundle        *AppBundleService
}

type service struct {
	c *Client
}

func NewCMAV2(config client.ClientConfig) (cma_service.SpaceIdClientBuilder, error) {
	return cma.New(config)
}

func NewCDAV2(config client.ClientConfig) (cda_service.SpaceIdClientBuilder, error) {
	return cda.New(config)
}

// NewCMA returns a CMA client
func NewCMA(token string) *Client {
	c := &Client{
		client: http.DefaultClient,
		api:    "CMA",
		token:  token,
		Debug:  false,
		Headers: map[string]string{
			"Authorization":           fmt.Sprintf("Bearer %s", token),
			"Content-Type":            "application/vnd.contentful.management.v1+json",
			"X-Contentful-User-Agent": fmt.Sprintf("sdk contentful.go/%s", Version),
		},
		BaseURL:     "https://api.contentful.com",
		Environment: "master",
	}
	c.commonService.c = c

	c.Spaces = (*SpacesService)(&c.commonService)
	c.Users = (*UsersService)(&c.commonService)
	c.Organizations = (*OrganizationsService)(&c.commonService)
	c.Roles = (*RolesService)(&c.commonService)
	c.Memberships = (*MembershipsService)(&c.commonService)
	c.Snapshots = (*SnapshotsService)(&c.commonService)
	c.AccessTokens = (*AccessTokensService)(&c.commonService)
	c.EntryTasks = (*EntryTasksService)(&c.commonService)
	c.ScheduledActions = (*ScheduledActionsService)(&c.commonService)
	c.Webhooks = (*WebhooksService)(&c.commonService)
	c.WebhookCalls = (*WebhookCallsService)(&c.commonService)
	c.EditorInterfaces = (*EditorInterfacesService)(&c.commonService)
	c.Extensions = (*ExtensionsService)(&c.commonService)
	c.Usages = (*UsagesService)(&c.commonService)
	c.AppUpload = &AppUploadService{
		service: c.commonService,
		BaseURL: "https://upload.contentful.com",
	}

	c.AppBundle = (*AppBundleService)(&c.commonService)
	return c
}

// NewCDA returns a CDA client
func NewCDA(token string) *Client {
	c := &Client{
		client: http.DefaultClient,
		api:    "CDA",
		token:  token,
		Debug:  false,
		Headers: map[string]string{
			"Authorization":           "Bearer " + token,
			"Content-Type":            "application/vnd.contentful.delivery.v1+json",
			"X-Contentful-User-Agent": fmt.Sprintf("contentful-go/%s", Version),
		},
		BaseURL:     "https://cdn.contentful.com",
		Environment: "master",
	}
	c.commonService.c = c

	c.Spaces = (*SpacesService)(&c.commonService)
	c.Webhooks = (*WebhooksService)(&c.commonService)

	return c
}

// NewCPA returns a CPA client
func NewCPA(token string) *Client {
	c := &Client{
		client: http.DefaultClient,
		Debug:  false,
		api:    "CPA",
		token:  token,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BaseURL: "https://preview.contentful.com",
	}

	c.Spaces = &SpacesService{c: c}
	c.Webhooks = &WebhooksService{c: c}

	return c
}

// NewResourceClient returns a client for the resource/uploads endpoints
func NewResourceClient(token string) *Client {
	c := &Client{
		client: http.DefaultClient,
		api:    "URC",
		Debug:  false,
		token:  token,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BaseURL: "https://upload.contentful.com",
	}
	c.commonService.c = c

	c.Resources = (*ResourcesService)(&c.commonService)

	return c
}

// SetOrganization sets the given organization id
func (c *Client) SetOrganization(organizationID string) *Client {
	c.Headers["X-Contentful-Organization"] = organizationID

	return c
}

// SetEnvironment sets the given environment.
// https://www.contentful.com/developers/docs/references/content-management-api/#/reference/environments
func (c *Client) SetEnvironment(environment string) *Client {
	c.Environment = environment
	return c
}

// SetHTTPClient sets the underlying http.Client used to make requests.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.client = client
}

func (c *Client) newRequest(method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	return c.newRequestWithBaseUrl(method, c.BaseURL, path, query, body)
}

func (c *Client) newRequestWithBaseUrl(method, baseUrl string, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	// set query params
	for key, value := range c.QueryParams {
		query.Set(key, value)
	}

	u.Path = path
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	if c.Debug {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 200 && res.StatusCode < 400 {
		// Upload/Create Resource response cannot be decoded
		if c.api == "URC" && req.Method == "POST" {
			defer res.Body.Close()
		} else {
			if v != nil {
				defer res.Body.Close()
				err = json.NewDecoder(res.Body).Decode(v)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// parse api response
	apiError := c.handleError(req, res)

	// return apiError if it is not rate limit error
	if _, ok := apiError.(common.RateLimitExceededError); !ok {
		return apiError
	}

	resetHeader := res.Header.Get("x-contentful-ratelimit-reset")

	// return apiError if Ratelimit-Reset header is not presented
	if resetHeader == "" {
		return apiError
	}

	// wait X-Contentful-Ratelimit-Reset amount of seconds
	waitSeconds, err := strconv.Atoi(resetHeader)
	if err != nil {
		return apiError
	}

	time.Sleep(time.Second * time.Duration(waitSeconds))

	return c.do(req, v)
}

func (c *Client) handleError(req *http.Request, res *http.Response) error {
	if c.Debug {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%q", dump)
	}

	var e common.ErrorResponse
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&e)
	if err != nil {
		return err
	}

	apiError := common.NewApiError(req, res, &e)

	switch errType := e.Sys.ID; errType {
	case "NotFound":
		return common.NotFoundError{apiError}
	case "RateLimitExceeded":
		return common.RateLimitExceededError{apiError}
	case "AccessTokenInvalid":
		return common.AccessTokenInvalidError{apiError}
	case "ValidationFailed":
		return common.ValidationFailedError{apiError}
	case "VersionMismatch":
		return common.VersionMismatchError{apiError}
	case "Conflict":
		return common.VersionMismatchError{apiError}
	case "InvalidEntry":
		return common.InvalidEntryError{apiError}
	default:
		return e
	}
}
