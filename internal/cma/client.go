package cma

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	common2 "github.com/flaconi/contentful-go/pkgs/common"
	"github.com/flaconi/contentful-go/service"
	"github.com/flaconi/contentful-go/service/common"
	"io"
	"log"
	"moul.io/http2curl"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

var _ common.SpaceIdClientBuilder = &Client{}

type Client struct {
	httpClient common.HttpClient
	debug      bool
	headers    map[string]string
	url        *url.URL
}

func New(config ClientConfig) (common.SpaceIdClientBuilder, error) {

	httpClient := config.HTTPClient

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	userAgent := config.UserAgent

	if userAgent == "" {
		userAgent = fmt.Sprintf("sdk contentful.go/%s", service.Version)
	}

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return nil, err
	}
	return &Client{
		httpClient: httpClient,
		debug:      config.Debug,
		headers: map[string]string{
			"Authorization":           fmt.Sprintf("Bearer %s", config.Token),
			"Content-Type":            "application/vnd.contentful.management.v1+json",
			"X-Contentful-User-Agent": userAgent,
		},
		url: parsedURL,
	}, nil
}

func (c *Client) WithSpaceId(spaceId string) common.SpaceIdClient {
	return &SpaceIdClient{
		client:  c,
		spaceId: spaceId,
	}
}

func (c *Client) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.execute(ctx, http.MethodGet, path, queryParams, headers, nil)
}

func (c *Client) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.execute(ctx, http.MethodPost, path, queryParams, headers, body)
}

func (c *Client) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.execute(ctx, http.MethodPut, path, queryParams, headers, body)
}

func (c *Client) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.execute(ctx, http.MethodDelete, path, queryParams, headers, nil)
}

func (c *Client) createEndpoint(p string) (*url.URL, error) {
	parsedURL, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	return c.url.ResolveReference(parsedURL), nil
}

func (c *Client) execute(ctx context.Context, method string, path string, params url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	endpoint, err := c.createEndpoint(path)
	if err != nil {
		return nil, err
	}

	if params != nil {
		endpoint.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("creating new request: %w", err)
	}

	if headers != nil {
		req.Header = headers
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	if c.debug {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 200 && res.StatusCode < 400 {
		return res, nil
	}

	// parse api response
	apiError := c.handleError(req, res)

	// return apiError if it is not rate limit error
	var rateLimitExceededError common2.RateLimitExceededError
	if !errors.As(apiError, &rateLimitExceededError) {
		return nil, apiError
	}

	resetHeader := res.Header.Get("x-contentful-ratelimit-reset")

	// return apiError if Ratelimit-Reset header is not presented
	if resetHeader == "" {
		return nil, apiError
	}

	// wait X-Contentful-Ratelimit-Reset amount of seconds
	waitSeconds, err := strconv.Atoi(resetHeader)
	if err != nil {
		return nil, apiError
	}

	time.Sleep(time.Second * time.Duration(waitSeconds))

	return c.execute(ctx, method, path, params, headers, body)
}

func (c *Client) handleError(req *http.Request, res *http.Response) error {
	if c.debug {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%q", dump)
	}

	var e common2.ErrorResponse
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&e)
	if err != nil {
		return err
	}

	apiError := common2.NewApiError(req, res, &e)

	switch errType := e.Sys.ID; errType {
	case "NotFound":
		return common2.NotFoundError{apiError}
	case "RateLimitExceeded":
		return common2.RateLimitExceededError{apiError}
	case "AccessTokenInvalid":
		return common2.AccessTokenInvalidError{apiError}
	case "ValidationFailed":
		return common2.ValidationFailedError{apiError}
	case "VersionMismatch":
		return common2.VersionMismatchError{apiError}
	case "Conflict":
		return common2.VersionMismatchError{apiError}
	default:
		return e
	}
}
