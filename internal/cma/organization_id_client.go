package cma

import (
	"context"
	"fmt"
	"github.com/flaconi/contentful-go/internal/cma/app_definitions"
	"github.com/flaconi/contentful-go/service/cma"
	"io"
	"net/http"
	"net/url"
)

var _ cma.OrganizationIdClient = &OrganizationIdClient{}

type OrganizationIdClient struct {
	client         *Client
	organizationId string
}

func (c *OrganizationIdClient) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Get(ctx, fmt.Sprintf("/organizations/%s%s", c.organizationId, path), queryParams, headers)
}

func (c *OrganizationIdClient) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Post(ctx, fmt.Sprintf("/organizations/%s%s", c.organizationId, path), queryParams, headers, body)
}

func (c *OrganizationIdClient) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Put(ctx, fmt.Sprintf("/organizations/%s%s", c.organizationId, path), queryParams, headers, body)
}

func (c *OrganizationIdClient) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Delete(ctx, fmt.Sprintf("/organizations/%s%s", c.organizationId, path), queryParams, headers)
}

func (c *OrganizationIdClient) AppDefinitions() cma.AppDefinitions {
	return app_definitions.NewAppDefinitionService(c)
}
