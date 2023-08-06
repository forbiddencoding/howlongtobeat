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
	"math"
)

type SearchGameSimple struct {
	GameID          int     `json:"game_id"`
	GameName        string  `json:"game_name"`
	ProfilePlatform string  `json:"profile_platform"`
	GameImage       string  `json:"game_image"`
	CompMain        float64 `json:"comp_main"`
	CompPlus        float64 `json:"comp_plus"`
	CompAll         float64 `json:"comp_all"`
	Similarity      float64 `json:"similarity"`
}

func (s *SearchGame) Reduce() []*SearchGameSimple {
	var reduced = make([]*SearchGameSimple, len(s.Data))

	for i, v := range s.Data {
		reduced[i] = &SearchGameSimple{
			GameID:          v.GameID,
			GameName:        v.GameName,
			ProfilePlatform: v.ProfilePlatform,
			GameImage:       v.GameImage,
			CompMain:        math.Round(float64(v.CompMain) / 3600),
			CompPlus:        math.Round(float64(v.CompPlus) / 3600),
			CompAll:         math.Round(float64(v.CompAll) / 3600),
			Similarity:      v.Similarity,
		}
	}

	return reduced
}

func (c *Client) SearchSimple(ctx context.Context, searchTerm string, searchModifier SearchModifier) ([]*SearchGameSimple, error) {
	result, err := c.Search(ctx, searchTerm, searchModifier, nil)
	if err != nil {
		return nil, err
	}

	return result.Reduce(), nil
}

type GameDetailSimple struct {
	GameID          int     `json:"game_id"`
	GameName        string  `json:"game_name"`
	ProfilePlatform string  `json:"profile_platform"`
	GameImage       string  `json:"game_image"`
	CompMain        float64 `json:"comp_main"`
	CompPlus        float64 `json:"comp_plus"`
	CompAll         float64 `json:"comp_all"`
}

func (s *GameDetails) Reduce() *GameDetailSimple {
	return &GameDetailSimple{
		GameID:          s.Props.PageProps.Game.Data.Game[0].GameID,
		GameName:        s.Props.PageProps.Game.Data.Game[0].GameName,
		ProfilePlatform: s.Props.PageProps.Game.Data.Game[0].ProfilePlatform,
		GameImage:       s.Props.PageProps.Game.Data.Game[0].GameImage,
		CompMain:        math.Round(float64(s.Props.PageProps.Game.Data.Game[0].CompMain) / 3600),
		CompPlus:        math.Round(float64(s.Props.PageProps.Game.Data.Game[0].CompPlus) / 3600),
		CompAll:         math.Round(float64(s.Props.PageProps.Game.Data.Game[0].CompAll) / 3600),
	}
}

func (c *Client) DetailSimple(ctx context.Context, gameID int) (*GameDetailSimple, error) {
	result, err := c.Detail(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return result.Reduce(), nil
}
