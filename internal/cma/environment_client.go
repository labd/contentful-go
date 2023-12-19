package cma

import (
	"context"
	"github.com/flaconi/contentful-go/service/common"
	"io"
	"net/http"
	"net/url"
)

var _ common.EnvironmentClient = &EnvironmentClient{}

type EnvironmentClient struct {
	Client      common.RestClient
	Environment string
}

func (c *EnvironmentClient) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *EnvironmentClient) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *EnvironmentClient) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (c *EnvironmentClient) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}
