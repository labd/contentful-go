package common

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type RestClient interface {
	Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error)
	Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error)
	Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error)
	Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error)
}

type SpaceIdClientBuilder interface {
	RestClient
	WithSpaceId(spaceId string) SpaceIdClient
}

type SpaceIdClient interface {
	RestClient
	WithEnvironment(environment string) EnvironmentClient
}

type EnvironmentClient interface {
	RestClient
}
