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
