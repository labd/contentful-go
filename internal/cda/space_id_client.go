package cda

import (
	"context"
	"fmt"
	"github.com/flaconi/contentful-go/service/cda"
	"io"
	"net/http"
	"net/url"
)

var _ cda.SpaceIdClient = &SpaceIdClient{}

type SpaceIdClient struct {
	client  *Client
	spaceId string
}

func (c *SpaceIdClient) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Get(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers)
}

func (c *SpaceIdClient) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Post(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers, body)
}

func (c *SpaceIdClient) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Put(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers, body)
}

func (c *SpaceIdClient) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Delete(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers)
}

func (c *SpaceIdClient) WithEnvironment(environment string) cda.EnvironmentClient {
	return &EnvironmentClient{
		client:      c,
		environment: environment,
	}
}
