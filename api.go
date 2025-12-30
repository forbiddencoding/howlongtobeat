package howlongtobeat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ApiData struct {
	token        string
	scriptPaths  []string
	endpointPath string
}

func (c *Client) setDefaultRequestHeaders(req *http.Request) {
	// Setting headers to match the browser request.
	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set(http.CanonicalHeaderKey("Origin"), "https://howlongtobeat.com/")
	req.Header.Set(http.CanonicalHeaderKey("Referer"), "https://howlongtobeat.com/")
}

func (c *Client) tokenHTTPRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		hltbTokenURL+"?t="+time.Now().Format(time.RFC3339Nano),
		nil,
	)
	if err != nil {
		return nil, err
	}

	c.setDefaultRequestHeaders(req)

	return req, nil
}

func (c *Client) scriptPathHTTPRequest(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		hltbBaseURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	c.setDefaultRequestHeaders(req)

	return req, nil
}

func (c *Client) endpointPathHTTPRequest(ctx context.Context, path string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		hltbBaseURL+path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	c.setDefaultRequestHeaders(req)

	return req, nil
}

// getApiDataWithEndpointSearch
// Method parses the request token and tries to search js scripts and then
// find the endpointPath in one of these scripts.
func (c *Client) getApiDataWithEndpointSearch(ctx context.Context) (*ApiData, error) {
	apiData := &ApiData{}

	req, err := c.tokenHTTPRequest(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create token request: %s.", err))
	}

	if err = c.do(req, c.tokenParser(apiData)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch token: %s.", err))
	}

	req, err = c.scriptPathHTTPRequest(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create script path request: %s.", err))
	}

	if err = c.do(req, c.scriptParser(apiData)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch script path: %s.", err))
	}

	for _, scriptPath := range apiData.scriptPaths {
		req, err = c.endpointPathHTTPRequest(ctx, scriptPath)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to create endpoint request: %s.", err))
		}

		if err = c.do(req, c.endpointParser(apiData)); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to fetch endpoint: %s.", err))
		}

		if apiData.endpointPath != "" {
			break
		}
	}

	if apiData.endpointPath == "" {
		return nil, errors.New("failed to find endpoint path")
	}

	return apiData, nil
}

// getApiDataWithDefaultEndpoint
// Method parses the request token and sets the default endpointPath.
func (c *Client) getApiDataWithDefaultEndpoint(ctx context.Context) (*ApiData, error) {
	apiData := &ApiData{}

	req, err := c.tokenHTTPRequest(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create token request: %s.", err))
	}

	if err = c.do(req, c.tokenParser(apiData)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch token: %s.", err))
	}

	apiData.endpointPath = hltbSearchEndpoint

	return apiData, nil
}

// getApiData
// Search flag indicates which getApi method we should use.
// If true, the method uses getApi method with endpoint search, otherwise
// it uses the default endpoint.
func (c *Client) getApiData(ctx context.Context, search bool) (*ApiData, error) {
	if c.apiData != nil {
		return c.apiData, nil
	}

	var (
		apiData *ApiData
		err     error
	)
	if search {
		apiData, err = c.getApiDataWithEndpointSearch(ctx)
	} else {
		apiData, err = c.getApiDataWithDefaultEndpoint(ctx)
	}

	if err != nil {
		return nil, err
	}

	c.apiData = apiData

	return c.apiData, nil
}
