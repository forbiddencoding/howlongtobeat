package howlongtobeat

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// hltbBaseURL is the base URL for the HowLongToBeat.
	hltbBaseURL = "https://howlongtobeat.com"
	// hltbSearchEndpoint is the default endpoint for the HowLongToBeat search API.
	hltbSearchEndpoint = "/api/finder"
	// hltbTokenURL is the URL to retrieve the token for the HowLongToBeat API.
	hltbTokenURL = "https://howlongtobeat.com/api/finder/init"
	// hltbGameURL is the base URL for the HowLongToBeat game API.
	hltbGameURL = "https://howlongtobeat.com/game"
	// defaultRequestTimeout is the default timeout for outgoing requests, we wait up to 30 seconds.
	defaultRequestTimeout = 30 * time.Second
)

type (
	Client struct {
		client  *http.Client
		logger  *log.Logger
		apiData *ApiData
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
		return errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}
