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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type (
	// SearchGameData contains completion levels and their respective time, as well as some other data.
	SearchGameData struct {
		Count           int    `json:"count"`
		GameID          int    `json:"game_id"`
		GameName        string `json:"game_name"`
		GameNameDate    int    `json:"game_name_date"`
		GameAlias       string `json:"game_alias"`
		GameType        string `json:"game_type"`
		GameImage       string `json:"game_image"`
		CompLvlCombine  int    `json:"comp_lvl_combine"`
		CompLvlSp       int    `json:"comp_lvl_sp"`
		CompLvlCo       int    `json:"comp_lvl_co"`
		CompLvlMp       int    `json:"comp_lvl_mp"`
		CompLvlSpd      int    `json:"comp_lvl_spd"`
		CompMain        int    `json:"comp_main"`
		CompPlus        int    `json:"comp_plus"`
		Comp100         int    `json:"comp_100"`
		CompAll         int    `json:"comp_all"`
		CompMainCount   int    `json:"comp_main_count"`
		CompPlusCount   int    `json:"comp_plus_count"`
		Comp100Count    int    `json:"comp_100_count"`
		CompAllCount    int    `json:"comp_all_count"`
		InvestedCo      int    `json:"invested_co"`
		InvestedMp      int    `json:"invested_mp"`
		InvestedCoCount int    `json:"invested_co_count"`
		InvestedMpCount int    `json:"invested_mp_count"`
		CountComp       int    `json:"count_comp"`
		CountSpeedrun   int    `json:"count_speedrun"`
		CountBacklog    int    `json:"count_backlog"`
		CountReview     int    `json:"count_review"`
		ReviewScore     int    `json:"review_score"`
		CountPLaying    int    `json:"count_playing"`
		CountRetired    int    `json:"count_retired"`
		ProfileDev      string `json:"profile_dev"`
		ProfilePopular  int    `json:"profile_popular"`
		ProfileSteam    int    `json:"profile_steam"`
		ProfilePlatform string `json:"profile_platform"`
		ReleaseWorld    int    `json:"release_world"`
		// Similarity is not part of the JSON response. It is a calculated value that represents the similarity between
		// the search query and the game name. It is only an approximation of the similarity between the two strings, to
		// aid in sorting the results, in case the search query has more than one result.
		Similarity float64 `json:"-"`
	}

	SearchGame struct {
		Color       string           `json:"color"`
		Title       string           `json:"title"`
		Category    string           `json:"category"`
		Count       int              `json:"count"`
		PageCurrent int              `json:"page_current"`
		PageTotal   int              `json:"page_total"`
		PageSize    int              `json:"page_size"`
		Data        []SearchGameData `json:"data"`
	}

	SearchGamePagination struct {
		Page     int
		PageSize int
	}

	searchRequestOptionsGamesRangeTime struct {
		Min int `json:"min"`
		Max int `json:"max"`
	}

	searchRequestOptionsGamesGameplay struct {
		Difficulty  string `json:"difficulty"`
		Flow        string `json:"flow"`
		Genre       string `json:"genre"`
		Perspective string `json:"perspective"`
	}

	searchRequestOptionsGames struct {
		UserID        int                                `json:"userId,omitempty"`
		Platform      string                             `json:"platform"`
		SortCategory  string                             `json:"sortCategory,omitempty"`
		RangeCategory string                             `json:"rangeCategory,omitempty"`
		RangeTime     searchRequestOptionsGamesRangeTime `json:"rangeTime"`
		Gameplay      searchRequestOptionsGamesGameplay  `json:"gameplay"`
		Modifier      SearchModifier                     `json:"modifier,omitempty"`
	}

	searchRequestOptionsUsers struct {
		SortCategory string `json:"sortCategory,omitempty"`
	}

	searchRequestOptions struct {
		Games      searchRequestOptionsGames `json:"games,omitempty"`
		Users      searchRequestOptionsUsers `json:"users,omitempty"`
		Filter     string                    `json:"filter,omitempty"`
		Sort       int                       `json:"sort,omitempty"`
		Randomizer int                       `json:"randomizer,omitempty"`
	}

	searchRequest struct {
		SearchType    string               `json:"searchType"`
		SearchTerms   []string             `json:"searchTerms"`
		SearchPage    int                  `json:"searchPage"`
		Size          int                  `json:"size"`
		SearchOptions searchRequestOptions `json:"searchOptions"`
	}
)

type SearchModifier string

const (
	SearchModifierNone    SearchModifier = ""
	SearchModifierOnlyDLC SearchModifier = "only_dlc"
	SearchModifierHideDLC SearchModifier = "hide_dlc"
)

func (sm SearchModifier) String() string {
	return string(sm)
}

func (c *Client) prepSearchRequest(searchTerm string, searchModifier SearchModifier, pagination *SearchGamePagination) *searchRequest {
	requestBody := &searchRequest{
		SearchOptions: searchRequestOptions{
			Games: searchRequestOptionsGames{
				SortCategory:  "popular",
				RangeCategory: "main",
				Gameplay:      searchRequestOptionsGamesGameplay{},
				RangeTime: searchRequestOptionsGamesRangeTime{
					Min: 0,
					Max: 0,
				},
				Modifier: searchModifier,
			},
		},
		SearchType: "games",
	}

	requestBody.SearchTerms = strings.Split(searchTerm, " ")

	if pagination != nil {
		requestBody.SearchPage = c.normalizePaginationValue(pagination.Page, 1)
		requestBody.Size = c.normalizePaginationValue(pagination.PageSize, 20)
	} else {
		requestBody.SearchPage = 1
		requestBody.Size = 20
	}

	return requestBody
}

// normalizePaginationValue will return the pagination value if it's greater than zero
// and default value otherwise.
func (c *Client) normalizePaginationValue(value, defaultVal int) int {
	if value < 1 {
		return defaultVal
	}

	return value
}

func (c *Client) searchHTTPRequest(ctx context.Context, body []byte, endpoint, token string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		hltbBaseURL+"/"+endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	// Setting headers to match the browser request.
	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set(http.CanonicalHeaderKey("Origin"), "https://howlongtobeat.com/")
	req.Header.Set(http.CanonicalHeaderKey("Referer"), "https://howlongtobeat.com/")
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")
	req.Header.Set(http.CanonicalHeaderKey("Accept"), "*/*")
	req.Header.Set(http.CanonicalHeaderKey("Accept-Language"), "en")
	req.Header.Set(http.CanonicalHeaderKey("Cache-Control"), "no-cache")
	req.Header.Set(http.CanonicalHeaderKey("Pragma"), "no-cache")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Ch-Ua"), `Not;A Brand;v="99", "Google Chrome";v="115", "Chromium";v="115"`)
	req.Header.Set(http.CanonicalHeaderKey("Sec-Ch-Ua-Mobile"), "?0")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Fetch-Mode"), "cors")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Fetch-Dest"), "empty")
	req.Header.Set(http.CanonicalHeaderKey("Dnt"), "1")
	req.Header.Set(http.CanonicalHeaderKey("x-auth-token"), token)

	return req, nil
}

// Search searches for games on HowLongToBeat.
// SearchTerm is typically the title of the game or DLC.
// SearchModifier can be used to filter the results by either excluding or including games and DLCs.
// Pagination is optional, but recommended. The default page size is 20.
func (c *Client) Search(ctx context.Context, searchTerm string, searchModifier SearchModifier, pagination *SearchGamePagination) (*SearchGame, error) {
	if searchTerm == "" {
		return nil, errors.New("search term cannot be empty")
	}

	apiData, err := c.getApiData(ctx)
	if err != nil {
		return nil, err
	}

	requestBody := c.prepSearchRequest(searchTerm, searchModifier, pagination)

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := c.searchHTTPRequest(ctx, body, apiData.endpointPath, apiData.token)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create search request: %s.", err))
	}

	var resp SearchGame

	if err = c.do(req, c.jsonParser(&resp)); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to search: %s.", err))
	}

	var searchResults = make([]SearchGameData, len(resp.Data))

	for i, result := range resp.Data {
		result.Similarity = calculateJaccardSimilarity(searchTerm, result.GameName)

		searchResults[i] = result
	}

	resp.Data = searchResults

	if len(resp.Data) > 0 {
		sort.Slice(resp.Data, func(i, j int) bool {
			return resp.Data[i].Similarity > resp.Data[j].Similarity
		})
	}

	return &resp, nil
}
