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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	// GameDetailsGameDataGame contains metadata about the game and the completion data.
	GameDetailsGameDataGame struct {
		GameID            int    `json:"game_id"`
		GameName          string `json:"game_name"`
		GameNameDate      int    `json:"game_name_date"`
		CountPlaying      int    `json:"count_playing"`
		CountBacklog      int    `json:"count_backlog"`
		CountReplay       int    `json:"count_replay"`
		CountCustom       int    `json:"count_custom"`
		CountComp         int    `json:"count_comp"`
		CountRetired      int    `json:"count_retired"`
		CountReview       int    `json:"count_review"`
		ReviewScore       int    `json:"review_score"`
		GameAlias         string `json:"game_alias"`
		GameImage         string `json:"game_image"`
		GameType          string `json:"game_type"`
		GameParent        int    `json:"game_parent"`
		ProfileSummary    string `json:"profile_summary"`
		ProfileDev        string `json:"profile_dev"`
		ProfilePub        string `json:"profile_pub"`
		ProfilePlatform   string `json:"profile_platform"`
		ProfileGenre      string `json:"profile_genre"`
		ProfileSteam      int    `json:"profile_steam"`
		ProfileSteamAlt   int    `json:"profile_steam_alt"`
		ProfileItch       int    `json:"profile_itch"`
		ProfileIgn        string `json:"profile_ign"`
		ReleaseWorld      string `json:"release_world"`
		ReleaseNa         string `json:"release_na"`
		ReleaseEu         string `json:"release_eu"`
		ReleaseJp         string `json:"release_jp"`
		RatingEsrb        string `json:"rating_esrb"`
		RatingPegi        string `json:"rating_pegi"`
		RatingCero        string `json:"rating_cero"`
		CompLvlSp         int    `json:"comp_lvl_sp"`
		CompLvlSpd        int    `json:"comp_lvl_spd"`
		CompLvlCo         int    `json:"comp_lvl_co"`
		CompLvlMp         int    `json:"comp_lvl_mp"`
		CompLvlCombine    int    `json:"comp_lvl_combine"`
		CompLvlPlatform   int    `json:"comp_lvl_platform"`
		CompAllCount      int    `json:"comp_all_count"`
		CompAll           int    `json:"comp_all"`
		CompAllL          int    `json:"comp_all_l"`
		CompAllH          int    `json:"comp_all_h"`
		CompAllAvg        int    `json:"comp_all_avg"`
		CompAllMed        int    `json:"comp_all_med"`
		CompMainCount     int    `json:"comp_main_count"`
		CompMain          int    `json:"comp_main"`
		CompMainL         int    `json:"comp_main_l"`
		CompMainH         int    `json:"comp_main_h"`
		CompMainAvg       int    `json:"comp_main_avg"`
		CompMainMed       int    `json:"comp_main_med"`
		CompPlusCount     int    `json:"comp_plus_count"`
		CompPlus          int    `json:"comp_plus"`
		CompPlusL         int    `json:"comp_plus_l"`
		CompPlusH         int    `json:"comp_plus_h"`
		CompPlusAvg       int    `json:"comp_plus_avg"`
		CompPlusMed       int    `json:"comp_plus_med"`
		Comp100Count      int    `json:"comp_100_count"`
		Comp100           int    `json:"comp_100"`
		Comp100L          int    `json:"comp_100_l"`
		Comp100H          int    `json:"comp_100_h"`
		Comp100Avg        int    `json:"comp_100_avg"`
		Comp100Med        int    `json:"comp_100_med"`
		CompSpeedCount    int    `json:"comp_speed_count"`
		CompSpeed         int    `json:"comp_speed"`
		CompSpeedMin      int    `json:"comp_speed_min"`
		CompSpeedMax      int    `json:"comp_speed_max"`
		CompSpeedAvg      int    `json:"comp_speed_avg"`
		CompSpeedMed      int    `json:"comp_speed_med"`
		CompSpeed100Count int    `json:"comp_speed100_count"`
		CompSpeed100      int    `json:"comp_speed100"`
		CompSpeed100Min   int    `json:"comp_speed100_min"`
		CompSpeed100Max   int    `json:"comp_speed100_max"`
		CompSpeed100Avg   int    `json:"comp_speed100_avg"`
		CompSpeed100Med   int    `json:"comp_speed100_med"`
		CountTotal        int    `json:"count_total"`
		InvestedCoCount   int    `json:"invested_co_count"`
		InvestedCo        int    `json:"invested_co"`
		InvestedCoL       int    `json:"invested_co_l"`
		InvestedCoH       int    `json:"invested_co_h"`
		InvestedCoAvg     int    `json:"invested_co_avg"`
		InvestedCoMed     int    `json:"invested_co_med"`
		InvestedMpCount   int    `json:"invested_mp_count"`
		InvestedMp        int    `json:"invested_mp"`
		InvestedMpL       int    `json:"invested_mp_l"`
		InvestedMpH       int    `json:"invested_mp_h"`
		InvestedMpAvg     int    `json:"invested_mp_avg"`
		InvestedMpMed     int    `json:"invested_mp_med"`
		AddedStats        string `json:"added_stats"`
	}

	// GameDetailsGameDataIndividuality contains different completion types and times for a game based on a special
	// platform version e.g., PC release, 8th gen console release, etc.
	GameDetailsGameDataIndividuality struct {
		Platform  string `json:"platform"`
		CountComp string `json:"count_comp"`
		CompMain  string `json:"comp_main"`
		CompPlus  string `json:"comp_plus"`
		Comp100   string `json:"comp_100"`
		CompAll   string `json:"comp_all"`
		Compare   string `json:"compare"`
	}

	// GameDetailsGameDataRelationships contains information about the game's related content like DLCs.
	GameDetailsGameDataRelationships struct {
		GameID       int    `json:"game_id"`
		GameName     string `json:"game_name"`
		GameType     string `json:"game_type"`
		CompMain     int    `json:"comp_main"`
		CompPlus     int    `json:"comp_plus"`
		Comp100      int    `json:"comp_100"`
		CompAll      int    `json:"comp_all"`
		CompAllCount int    `json:"comp_all_count"`
		CountBacklog int    `json:"count_backlog"`
		ReviewScore  int    `json:"review_score"`
	}

	// GameDetailsGameDataUserReview contains user review scores as a count of reviews for each score
	// in an increment of five (5).
	GameDetailsGameDataUserReview struct {
		ReviewCount5   string `json:"5,omitempty"`
		ReviewCount10  string `json:"10,omitempty"`
		ReviewCount15  string `json:"15,omitempty"`
		ReviewCount20  string `json:"20,omitempty"`
		ReviewCount25  string `json:"25,omitempty"`
		ReviewCount30  string `json:"30,omitempty"`
		ReviewCount35  string `json:"35,omitempty"`
		ReviewCount40  string `json:"40,omitempty"`
		ReviewCount45  string `json:"45,omitempty"`
		ReviewCount50  string `json:"50,omitempty"`
		ReviewCount55  string `json:"55,omitempty"`
		ReviewCount60  string `json:"60,omitempty"`
		ReviewCount65  string `json:"65,omitempty"`
		ReviewCount70  string `json:"70,omitempty"`
		ReviewCount75  string `json:"75,omitempty"`
		ReviewCount80  string `json:"80,omitempty"`
		ReviewCount85  string `json:"85,omitempty"`
		ReviewCount90  string `json:"90,omitempty"`
		ReviewCount95  string `json:"95,omitempty"`
		ReviewCount100 string `json:"100,omitempty"`
		// ReviewCount is the total sum of all reviews.
		ReviewCount int `json:"review_count"`
	}

	// GameDetailsGameDataPlatformData contains information about the game's completion status on a specific platform
	// e.g., PC, PS4, Xbox One, etc.
	GameDetailsGameDataPlatformData struct {
		Platform   string `json:"platform"`
		CountComp  int    `json:"count_comp"`
		CountTotal int    `json:"count_total"`
		CompMain   int    `json:"comp_main"`
		CompPlus   int    `json:"comp_plus"`
		Comp100    int    `json:"comp_100"`
		CompLow    int    `json:"comp_low"`
		CompHigh   int    `json:"comp_high"`
	}

	GameDetailsIgnMap struct {
		Typename    string   `json:"__typename"`
		MapType     string   `json:"mapType"`
		MapName     string   `json:"mapName"`
		MapSlug     string   `json:"mapSlug"`
		ObjectSlug  string   `json:"objectSlug"`
		ObjectName  string   `json:"objectName"`
		Width       float64  `json:"width"`
		Height      float64  `json:"height"`
		MinZoom     int      `json:"minZoom"`
		MaxZoom     int      `json:"maxZoom"`
		Tilesets    []string `json:"tilesets"`
		InitialLat  float64  `json:"initialLat"`
		InitialLng  float64  `json:"initialLng"`
		InitialZoom int      `json:"initialZoom"`
	}

	// GameDetailsIgnWikiNav contains data about IGN uri that can be appended to the base url
	// `https://ign.com/wikis/{game-name}/{URL}`
	GameDetailsIgnWikiNav struct {
		Typename string `json:"__typename"`
		Label    string `json:"label"`
		URL      string `json:"url,omitempty"`
	}

	// GameDetailsPageMetadata contains metadata about the howlongtobeat.com page.
	// This data can usually be ignored.
	GameDetailsPageMetadata struct {
		Title       string `json:"title"`
		Image       string `json:"image"`
		Description string `json:"description"`
		Canonical   string `json:"canonical"`
		Template    string `json:"template"`
	}

	// GameDetailsQuery contains the internal howlongtobeat.com game id.
	GameDetailsQuery struct {
		GameID string `json:"gameId"`
	}

	GameDetailsGameData struct {
		Game          []GameDetailsGameDataGame          `json:"game"`
		Individuality []GameDetailsGameDataIndividuality `json:"individuality"`
		Relationships []GameDetailsGameDataRelationships `json:"relationships"`
		UserReviews   GameDetailsGameDataUserReview      `json:"userReviews"`
		PlatformData  []GameDetailsGameDataPlatformData  `json:"platformData"`
	}

	GameDetails struct {
		Props struct {
			PageProps struct {
				Game struct {
					Count int
					Data  GameDetailsGameData
				}
				IgnWikiSlug  string
				IgnMap       *GameDetailsIgnMap
				IgnWikiNav   []GameDetailsIgnWikiNav
				PageMetadata GameDetailsPageMetadata
			}
		}
		Page  string
		Query GameDetailsQuery
	}

	gameDetailsResponse struct {
		Props struct {
			PageProps struct {
				Game struct {
					Count int                 `json:"count"`
					Data  GameDetailsGameData `json:"data"`
				} `json:"game"`
				IgnWikiSlug  string                  `json:"ignWikiSlug"`
				IgnMap       *GameDetailsIgnMap      `json:"ignMap"`
				IgnWikiNav   json.RawMessage         `json:"ignWikiNav,omitempty"`
				PageMetadata GameDetailsPageMetadata `json:"pageMetadata"`
			} `json:"pageProps"`
		} `json:"props"`
		Page  string           `json:"page"`
		Query GameDetailsQuery `json:"query"`
	}
)

func (c *Client) detailHTTPRequest(gameID int) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", hltbGameURL, gameID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set(http.CanonicalHeaderKey("Accept"), "*/*")
	req.Header.Set(http.CanonicalHeaderKey("Accept-Language"), "en")
	req.Header.Set(http.CanonicalHeaderKey("Cache-Control"), "no-cache")
	req.Header.Set(http.CanonicalHeaderKey("Pragma"), "no-cache")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Ch-Ua"), `Not;A Brand;v="99", "Google Chrome";v="115", "Chromium";v="115"`)
	req.Header.Set(http.CanonicalHeaderKey("Sec-Ch-Ua-Mobile"), "?0")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Fetch-Mode"), "cors")
	req.Header.Set(http.CanonicalHeaderKey("Sec-Fetch-Dest"), "empty")
	req.Header.Set(http.CanonicalHeaderKey("Dnt"), "1")

	return req, nil
}

// Detail returns the details of a game by its HLTB ID.
// If the context expires, the request will be canceled.
// If the gameID is 0, an error will be returned.
func (c *Client) Detail(ctx context.Context, gameID int) (*GameDetails, error) {
	if gameID == 0 {
		return nil, errors.New("gameID is required")
	}

	req, err := c.detailHTTPRequest(gameID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create game details request: %s.", err))
	}

	var response gameDetailsResponse

	if err = c.do(ctx, req, c.htmlParserByID(&response, "__NEXT_DATA__")); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to execute game details request: %s.", err))
	}

	return response.convertResponseToGameDetails()
}

func (g *gameDetailsResponse) convertResponseToGameDetails() (*GameDetails, error) {
	var data GameDetails
	// copy the fields from the response to the result
	data.Props.PageProps.Game.Count = g.Props.PageProps.Game.Count
	data.Props.PageProps.Game.Data = g.Props.PageProps.Game.Data
	data.Props.PageProps.IgnWikiSlug = g.Props.PageProps.IgnWikiSlug
	data.Props.PageProps.PageMetadata = g.Props.PageProps.PageMetadata
	data.Page = g.Page
	data.Query = g.Query

	// Handle the special case IgnWikiNav field
	var tempStruct GameDetailsIgnWikiNav
	if err := json.Unmarshal(g.Props.PageProps.IgnWikiNav, &tempStruct); err != nil {
		var tempSlice []GameDetailsIgnWikiNav
		if err = json.Unmarshal(g.Props.PageProps.IgnWikiNav, &tempSlice); err != nil {
			return &data, errors.New("invalid IgnWikiNav format")
		}
		data.Props.PageProps.IgnWikiNav = tempSlice
	} else {
		// Only add to slice if not the zero value
		if (tempStruct != GameDetailsIgnWikiNav{}) {
			data.Props.PageProps.IgnWikiNav = []GameDetailsIgnWikiNav{tempStruct}
		}
	}

	if g.Props.PageProps.IgnMap != nil {
		data.Props.PageProps.IgnMap = g.Props.PageProps.IgnMap
	}

	return &data, nil
}
