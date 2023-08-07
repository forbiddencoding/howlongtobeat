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
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_normalizePaginationValue(t *testing.T) {
	tests := []struct {
		name       string
		value      int
		defaultVal int
		want       int
	}{
		{
			name:       "Test with negative value",
			value:      -1,
			defaultVal: 10,
			want:       10,
		},
		{
			name:       "Test with zero value",
			value:      0,
			defaultVal: 20,
			want:       20,
		},
		{
			name:       "Test with positive value",
			value:      15,
			defaultVal: 30,
			want:       15,
		},
	}

	mockClient := &Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockClient.normalizePaginationValue(tt.value, tt.defaultVal); got != tt.want {
				t.Errorf("normalizePaginationValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchModifierString(t *testing.T) {
	tests := []struct {
		name string
		sm   SearchModifier
		want string
	}{
		{
			name: "Test SearchModifierNone",
			sm:   SearchModifierNone,
			want: "",
		},
		{
			name: "Test SearchModifierOnlyDLC",
			sm:   SearchModifierOnlyDLC,
			want: "only_dlc",
		},
		{
			name: "Test SearchModifierHideDLC",
			sm:   SearchModifierHideDLC,
			want: "hide_dlc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.String(); got != tt.want {
				t.Errorf("SearchModifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prepSearchRequest(t *testing.T) {
	tests := []struct {
		name           string
		term           string
		pagination     *SearchGamePagination
		searchModifier SearchModifier
		want           *searchRequest
	}{
		{
			name:           "Test with empty pagination",
			term:           "test game",
			pagination:     nil,
			searchModifier: SearchModifierNone,
			want: &searchRequest{
				SearchOptions: searchRequestOptions{
					Games: searchRequestOptionsGames{
						SortCategory:  "popular",
						RangeCategory: "main",
						RangeTime: searchRequestOptionsGamesRangeTime{
							Min: 0,
							Max: 0,
						},
					},
				},
				SearchType:  "games",
				SearchTerms: []string{"test", "game"},
				SearchPage:  1,
				Size:        20,
			},
		},
		{
			name: "Test with provided pagination",
			term: "another game",
			pagination: &SearchGamePagination{
				Page:     2,
				PageSize: 50,
			},
			searchModifier: SearchModifierNone,
			want: &searchRequest{
				SearchOptions: searchRequestOptions{
					Games: searchRequestOptionsGames{
						SortCategory:  "popular",
						RangeCategory: "main",
						RangeTime: searchRequestOptionsGamesRangeTime{
							Min: 0,
							Max: 0,
						},
					},
				},
				SearchType:  "games",
				SearchTerms: []string{"another", "game"},
				SearchPage:  2,
				Size:        50,
			},
		},
	}

	mockClient := &Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mockClient.prepSearchRequest(tt.term, tt.searchModifier, tt.pagination); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepSearchRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchHTTPRequest(t *testing.T) {
	var (
		headers = map[string]string{
			http.CanonicalHeaderKey("User-Agent"):       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
			http.CanonicalHeaderKey("Origin"):           "https://howlongtobeat.com/",
			http.CanonicalHeaderKey("Referer"):          "https://howlongtobeat.com/",
			http.CanonicalHeaderKey("Content-Type"):     "application/json",
			http.CanonicalHeaderKey("Accept"):           "*/*",
			http.CanonicalHeaderKey("Accept-Language"):  "en",
			http.CanonicalHeaderKey("Cache-Control"):    "no-cache",
			http.CanonicalHeaderKey("Pragma"):           "no-cache",
			http.CanonicalHeaderKey("Sec-Ch-Ua"):        `Not;A Brand;v="99", "Google Chrome";v="115", "Chromium";v="115"`,
			http.CanonicalHeaderKey("Sec-Ch-Ua-Mobile"): "?0",
			http.CanonicalHeaderKey("Sec-Fetch-Mode"):   "cors",
			http.CanonicalHeaderKey("Sec-Fetch-Dest"):   "empty",
			http.CanonicalHeaderKey("Dnt"):              "1",
		}
		body       []byte
		mockClient = &Client{}
	)

	req, err := mockClient.searchHTTPRequest(body)
	if err != nil {
		t.Errorf("searchHTTPRequest() error = %v", err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("searchHTTPRequest() method = %v, want %v", req.Method, http.MethodPost)
	}

	if req.URL.String() != hltbSearchURL {
		t.Errorf("searchHTTPRequest() url = %v, want %v", req.URL.String(), hltbSearchURL)
	}

	var errs error

	for header, value := range headers {
		if req.Header.Get(header) != value {
			errs = errors.Join(errs, errors.New(fmt.Sprintf("detailHTTPRequest() did not set the correct %s header: want: %s, received: %s", header, value, req.Header.Get(header))))
		}
	}

	if errs != nil {
		t.Fatal(errs)
	}
}

func Test_Search(t *testing.T) {
	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := mockClient.Search(context.TODO(), "The Witcher 3 Wild Hunt", SearchModifierNone, nil)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Data[0].GameID != 10270 {
		t.Errorf("Search() gameID = %v, want %v", result.Data[0].GameID, 10270)
	}
}

func Test_Search_NullSearchTerm(t *testing.T) {
	mockClient := &Client{}

	_, err := mockClient.Search(context.TODO(), "", SearchModifierNone, nil)
	if strings.Compare(err.Error(), "search term cannot be empty") != 0 {
		t.Fatalf(`Search() expected "search term cannot be empty" error, but received: %v`, err)
	}
}
