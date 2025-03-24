package cda

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labd/contentful-go/internal/cda/sync"
	"github.com/labd/contentful-go/service/cda"
	"github.com/labd/contentful-go/service/common"
)

var _ cda.EnvironmentClient = &EnvironmentClient{}

type EnvironmentClient struct {
	client      common.RestClient
	environment string
}

func (c *EnvironmentClient) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Get(ctx, fmt.Sprintf("/environments/%s%s", c.environment, path), queryParams, headers)
}

func (c *EnvironmentClient) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Post(ctx, fmt.Sprintf("/environments/%s%s", c.environment, path), queryParams, headers, body)
}

func (c *EnvironmentClient) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Put(ctx, fmt.Sprintf("/environments/%s%s", c.environment, path), queryParams, headers, body)
}

func (c *EnvironmentClient) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Delete(ctx, fmt.Sprintf("/environments/%s%s", c.environment, path), queryParams, headers)
}

func (c *EnvironmentClient) Sync() cda.Sync {
	return sync.NewSyncService(c)
}
