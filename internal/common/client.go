package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	common2 "github.com/labd/contentful-go/pkgs/common"
	"github.com/labd/contentful-go/service/common"
	"moul.io/http2curl"
)

type ClientConfig struct {
	URL         *url.URL
	HTTPClient  common.HttpClient
	Debug       bool
	UserAgent   string
	ContentType string
	Token       string
	Logger      *slog.Logger
}

type Client struct {
	httpClient common.HttpClient
	debug      bool
	headers    map[string]string
	url        *url.URL
	logger     *slog.Logger
}

func NewInternalClient(config ClientConfig) *Client {
	httpClient := config.HTTPClient

	logger := config.Logger

	return &Client{
		logger:     logger,
		httpClient: httpClient,
		debug:      config.Debug,
		headers: map[string]string{
			"Authorization":           fmt.Sprintf("Bearer %s", config.Token),
			"Content-Type":            config.ContentType,
			"X-Contentful-User-Agent": config.UserAgent,
		},
		url: config.URL,
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

	var intermediateBody *bytes.Buffer

	if body != nil {
		intermediateBody = new(bytes.Buffer)
		_, err = intermediateBody.ReadFrom(body)

		if err != nil {
			c.logger.ErrorContext(ctx, fmt.Sprintf("Error reading body: %v", err))
			return nil, err
		}

		body = io.NopCloser(bytes.NewBuffer(intermediateBody.Bytes()))
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
		c.logger.DebugContext(ctx, command.String())
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

	if body != nil {
		body = io.NopCloser(bytes.NewBuffer(intermediateBody.Bytes()))
	}

	return c.execute(ctx, method, path, params, headers, body)
}

func (c *Client) handleError(req *http.Request, res *http.Response) error {
	//https://www.contentful.com/developers/docs/references/errors/
	if c.debug {
		dump, err := httputil.DumpResponse(res, true)
		if err != nil {
			log.Fatal(err)
		}

		dumpReq, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Fatal(err)
		}

		c.logger.Debug(fmt.Sprintf("%q", dump))
		c.logger.Debug(fmt.Sprintf("%q", dumpReq))
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
	case "InvalidEntry":
		return common2.InvalidEntryError{apiError}
	default:
		return e
	}
}
