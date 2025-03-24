package cma

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labd/contentful-go/internal/cma/app_installations"
	"github.com/labd/contentful-go/internal/cma/assets"
	"github.com/labd/contentful-go/internal/cma/content_types"
	"github.com/labd/contentful-go/internal/cma/entries"
	"github.com/labd/contentful-go/internal/cma/locales"
	"github.com/labd/contentful-go/service/cma"
	"github.com/labd/contentful-go/service/common"
)

var _ cma.EnvironmentClient = &EnvironmentClient{}

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

func (c *EnvironmentClient) AppInstallations() cma.AppInstallations {
	return app_installations.NewAppInstallationsService(c)
}

func (c *EnvironmentClient) Entries() cma.Entries {
	return entries.NewEntriesService(c)
}

func (c *EnvironmentClient) Assets() cma.Assets {
	return assets.NewAssetService(c)
}
func (c *EnvironmentClient) ContentTypes() cma.ContentTypes {
	return content_types.NewContentTypeService(c)
}

func (c *EnvironmentClient) Locales() cma.Locales {
	return locales.NewLocaleService(c)
}
