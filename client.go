package howlongtobeat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	Client struct {
		client  *http.Client
		logger  *log.Logger
		apiData *ApiData
	}

	// ApiData contains the data needed to make requests to the HLTB API.
	ApiData struct {
		token        string
		scriptPaths  []string
		endpointPath string
	}

	// Option is a type alias for functions to configure your Client.
	Option func(client *Client)
)

// WithRequestTimeout sets the timeout for outgoing requests.
// If timeout duration is set to 0, the default timeout of 30 seconds will be used.
// If using the WithHTTPClient option, make sure to set your client before the timeout.
func WithRequestTimeout(timeout int) Option {
	return func(client *Client) {
		if timeout > 0 {
			client.client.Timeout = time.Duration(timeout) * time.Second
		}
	}
}

// WithHTTPClient sets the user provided HTTP client to use it for outgoing requests.
// When using the WithRequestTimeout option, make sure to set your client before the timeout.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.client = httpClient
	}
}

// New creates a new HowLongToBeat client for optimized HTTP requests.
func New(options ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          10,
				IdleConnTimeout:       15 * time.Second,
				ResponseHeaderTimeout: 15 * time.Second,
				DisableKeepAlives:     false,
				ForceAttemptHTTP2:     true,
			},
			Timeout: defaultRequestTimeout,
		},
	}

	// Apply options
	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

// do performs the given request and parses the response with the provided parser.
func (c *Client) do(req *http.Request, parser parseResponseFunc) (err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return parser(resp)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

// request creates a new HTTP request with the default headers and context.
func (c *Client) request(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set(http.CanonicalHeaderKey("Origin"), "https://howlongtobeat.com/")
	req.Header.Set(http.CanonicalHeaderKey("Referer"), "https://howlongtobeat.com/")

	return req, nil
}

func (c *Client) getApiData(ctx context.Context) (*ApiData, error) {
	if c.apiData != nil {
		return c.apiData, nil
	}

	apiData, err := c.getApiDataWithDefaultEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	c.apiData = apiData

	return c.apiData, nil
}

// getApiDataWithDefaultEndpoint
// Method parses the request token and sets the default endpointPath.
func (c *Client) getApiDataWithDefaultEndpoint(ctx context.Context) (*ApiData, error) {
	apiData := &ApiData{}

	req, err := c.tokenHTTPRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}

	if err = c.do(req, c.tokenParser(apiData)); err != nil {
		return nil, fmt.Errorf("fetch token: %w", err)
	}

	apiData.endpointPath = hltbSearchEndpoint

	return apiData, nil
}

// getApiDataWithEndpointSearch
// Method parses the request token and tries to search js scripts and then
// find the endpointPath in one of these scripts.
func (c *Client) getApiDataWithEndpointSearch(ctx context.Context) (*ApiData, error) {
	apiData := &ApiData{}

	req, err := c.tokenHTTPRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}

	if err = c.do(req, c.tokenParser(apiData)); err != nil {
		return nil, fmt.Errorf("fetch token: %w", err)
	}

	req, err = c.scriptPathHTTPRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("create script path request: %w", err)
	}

	if err = c.do(req, c.scriptParser(apiData)); err != nil {
		return nil, fmt.Errorf("fetch script path: %w", err)
	}

	for _, scriptPath := range apiData.scriptPaths {
		req, err = c.endpointPathHTTPRequest(ctx, scriptPath)
		if err != nil {
			return nil, fmt.Errorf("create endpoint request: %w", err)
		}

		if err = c.do(req, c.endpointParser(apiData)); err != nil {
			return nil, fmt.Errorf("fetch endpoint: %w", err)
		}

		if apiData.endpointPath != "" {
			break
		}
	}

	if apiData.endpointPath == "" {
		return nil, errors.New("empty endpoint path")
	}

	return apiData, nil
}

func (c *Client) tokenHTTPRequest(ctx context.Context) (*http.Request, error) {
	req, err := c.request(ctx, http.MethodGet, hltbTokenURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))
	req.URL.RawQuery = q.Encode()

	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=utf-8")

	return req, nil
}

func (c *Client) scriptPathHTTPRequest(ctx context.Context) (*http.Request, error) {
	return c.request(ctx, http.MethodGet, hltbBaseURL, nil)
}

func (c *Client) endpointPathHTTPRequest(ctx context.Context, path string) (*http.Request, error) {
	return c.request(ctx, http.MethodGet, hltbBaseURL+path, nil)
}
