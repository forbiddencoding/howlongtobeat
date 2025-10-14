/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2023. Edgar Schmidt
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the
 * following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following
 * disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following
 * disclaimer in the documentation and/or other materials provided with the distribution.
 *
 * Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products
 * derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
 * INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
	// hltbSearchURL is the base URL for the HowLongToBeat search API.
	hltbSearchURL = "https://howlongtobeat.com/api/locate"
	// hltbGameURL is the base URL for the HowLongToBeat game API.
	hltbGameURL = "https://howlongtobeat.com/game"
	// defaultRequestTimeout is the default timeout for outgoing requests, we wait up to 30 seconds.
	defaultRequestTimeout = 30 * time.Second
)

type (
	Client struct {
		client  *http.Client
		logger  *log.Logger
		timeout time.Duration
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
			client.timeout = time.Duration(timeout) * time.Second
		} else {
			client.timeout = defaultRequestTimeout
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
		// Default HTTP client
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
		timeout: defaultRequestTimeout,
		apiData: &ApiData{},
	}

	// Apply options
	for _, opt := range options {
		opt(c)
	}

	c.client.Timeout = c.timeout

	return c, nil
}

// do performs the given request and parses the response with the provided parser.
func (c *Client) do(req *http.Request, parser parseResponseFunc) (err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return parser(resp)
	default:
		return errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
}
