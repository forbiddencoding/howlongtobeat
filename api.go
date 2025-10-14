package howlongtobeat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type ApiData struct {
	scriptPath   string
	endpointPath string
}

func (c *Client) setDefaultRequestHeaders(req *http.Request) {
	// Setting headers to match the browser request.
	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set(http.CanonicalHeaderKey("Origin"), "https://howlongtobeat.com/")
	req.Header.Set(http.CanonicalHeaderKey("Referer"), "https://howlongtobeat.com/")
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

func (c *Client) endpointPathHTTPRequest(path string, ctx context.Context) (*http.Request, error) {
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

func (c *Client) getApiData() (*ApiData, error) {
	if c.apiData.endpointPath != "" {
		return c.apiData, nil
	}

	apiData := &ApiData{}

	req, err := c.scriptPathHTTPRequest(context.Background())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create script path request: %s.", err))
	}

	if err = c.do(req, c.scriptParser(apiData)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch script path: %s.", err))
	}

	req, err = c.endpointPathHTTPRequest(apiData.scriptPath, context.Background())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create endpoint request: %s.", err))
	}

	if err = c.do(req, c.endpointParser(apiData)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch endpoint: %s.", err))
	}

	c.apiData = apiData

	return c.apiData, nil
}
