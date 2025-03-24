package cma

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	internalcommon "github.com/labd/contentful-go/internal/common"
	"github.com/labd/contentful-go/pkgs/client"
	"github.com/labd/contentful-go/pkgs/util"
	"github.com/labd/contentful-go/service"
	"github.com/labd/contentful-go/service/cma"
)

var _ cma.SpaceIdClientBuilder = &Client{}

type Client struct {
	client *internalcommon.Client
}

func New(config client.ClientConfig) (cma.SpaceIdClientBuilder, error) {

	httpClient := config.HTTPClient

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	userAgent := config.UserAgent

	if userAgent == nil {
		userAgent = util.ToPointer(fmt.Sprintf("sdk contentful.go/%s", service.Version))
	}

	configUrl := config.URL

	if configUrl == nil {
		configUrl = util.ToPointer("https://api.contentful.com")
	}

	parsedURL, err := url.Parse(*configUrl)
	if err != nil {
		return nil, err
	}

	logger := config.Logger

	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return &Client{
		client: internalcommon.NewInternalClient(internalcommon.ClientConfig{
			URL:         parsedURL,
			HTTPClient:  httpClient,
			Debug:       false,
			UserAgent:   *userAgent,
			ContentType: "application/vnd.contentful.management.v1+json",
			Token:       config.Token,
			Logger:      logger,
		}),
	}, nil
}

func (c *Client) WithSpaceId(spaceId string) cma.SpaceIdClient {
	return &SpaceIdClient{
		client:  c,
		spaceId: spaceId,
	}
}

func (c *Client) WithOrganizationId(organizationId string) cma.OrganizationIdClient {
	return &OrganizationIdClient{
		client:         c,
		organizationId: organizationId,
	}
}

func (c *Client) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Get(ctx, path, queryParams, headers)
}

func (c *Client) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Post(ctx, path, queryParams, headers, body)
}

func (c *Client) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Put(ctx, path, queryParams, headers, body)
}

func (c *Client) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Delete(ctx, path, queryParams, headers)
}
