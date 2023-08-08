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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

type parseResponseFunc func(resp *http.Response) error

// jsonParser returns a function that will decode the body of an http.Response as JSON into the provided struct.
func (c *Client) jsonParser(val any) parseResponseFunc {
	return func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(val)
	}
}

// htmlParserByID returns a function that will parse the HTML body of an http.Response, and unmarshal the JSON found
// in the first child of the element with the provided ID into the provided struct.
func (c *Client) htmlParserByID(val any, ID string) parseResponseFunc {
	return func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return c.parseHTML(body, ID, val)
	}
}

func (c *Client) htmlScriptDataParserByID(val any, ID string) parseResponseFunc {
	return func(resp *http.Response) error {
		startTag := []byte(fmt.Sprintf(`<script id="%s" type="application/json">`, ID))
		endTag := []byte(`</script>`)

		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)

		var data []byte
		startTagFound := false

		for scanner.Scan() {
			line := scanner.Bytes()

			if !startTagFound {
				startIndex := bytes.Index(line, startTag)
				if startIndex != -1 {
					startTagFound = true
					// Adjust line to start at the beginning of JSON content, not at the start tag
					line = line[startIndex+len(startTag):]
				}
			}

			if startTagFound {
				endIndex := bytes.Index(line, endTag)
				if endIndex != -1 {
					// Adjust line to end at the beginning of the end tag, not include the end tag
					line = line[:endIndex]
				}

				data = append(data, line...)

				if endIndex != -1 {
					break
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if !startTagFound {
			return errors.New("start tag not found")
		}

		return json.Unmarshal(data, &val)
	}
}

// parseHTML reads the HTML from the provided bytes, finds the first element with the provided ID,
// and unmarshal the JSON found in its first child into the provided struct.
func (c *Client) parseHTML(body []byte, ID string, val any) error {
	document, err := html.Parse(io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return err
	}

	element := c.getElementByID(document, ID)
	if element == nil || element.Data == "" {
		return errors.New("element not found")
	}

	if element.FirstChild == nil || element.FirstChild.Data == "" {
		return errors.New("element first child not found")
	}

	return json.Unmarshal([]byte(element.FirstChild.Data), &val)
}

func (c *Client) getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func (c *Client) checkID(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := c.getAttribute(n, "id")

		if ok && s == id {
			return true
		}
	}

	return false
}

func (c *Client) traverseTree(n *html.Node, id string) *html.Node {
	if c.checkID(n, id) {
		return n
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		res := c.traverseTree(child, id)

		if res != nil {
			return res
		}
	}

	return nil
}

func (c *Client) getElementByID(n *html.Node, id string) *html.Node {
	return c.traverseTree(n, id)
}
