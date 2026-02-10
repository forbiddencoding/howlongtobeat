package howlongtobeat

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type parseResponseFunc func(resp *http.Response) error

// jsonParser returns a function that will decode the body of an http.Response as JSON into the provided struct.
func (c *Client) jsonParser(val any) parseResponseFunc {
	return func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(val)
	}
}

func (c *Client) nextDataParser(val any) parseResponseFunc {
	return func(resp *http.Response) error {
		startTag := []byte(`<script id="__NEXT_DATA__" type="application/json">`)
		endTag := []byte(`</script>`)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		start := bytes.Index(body, startTag)
		end := bytes.Index(body[start:], endTag)

		return json.Unmarshal(body[start+len(startTag):start+end], &val)
	}
}

func (c *Client) scriptParser(apiData *ApiData) parseResponseFunc {
	return func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		reg := regexp.MustCompile(`(?i)<script[^>]*src=["']([^"']*/_next/static/chunks/[^"']*\.js)["'][^>]*>`)
		matches := reg.FindAllSubmatch(body, -1)

		if len(matches) == 0 {
			return errors.New("script src path not found")
		}

		apiData.scriptPaths = make([]string, len(matches))
		for i, match := range matches {
			apiData.scriptPaths[i] = string(match[1])
		}

		return nil
	}
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (c *Client) tokenParser(apiData *ApiData) parseResponseFunc {
	return func(resp *http.Response) error {
		var tokenResponse TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
			return err
		}

		apiData.token = tokenResponse.Token

		return nil
	}
}

func (c *Client) endpointParser(apiData *ApiData) parseResponseFunc {
	return func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		reg := regexp.MustCompile(`(?si)fetch\s*\(\s*["']/api/([a-zA-Z0-9_/]+)[^"']*["']\s*,\s*{[^}]*method:\s*["']POST["'][^}]*}`)
		matches := reg.FindSubmatch(body)
		if len(matches) < 2 {
			return errors.New("endpoint path not found")
		}

		var basePath string
		pathSuffix := string(matches[1])
		if strings.Contains(pathSuffix, "/") {
			basePath = strings.Split(pathSuffix, "/")[0]
		} else {
			basePath = pathSuffix
		}

		apiData.endpointPath = "/api/" + basePath

		return nil
	}
}
