package howlongtobeat

import (
	"math"
	"reflect"
	"testing"
)

func Test_SearchGame_Reduce(t *testing.T) {
	data := []SearchGameData{
		{
			GameID:          1,
			GameName:        "Game1",
			ProfilePlatform: "Platform1",
			GameImage:       "Image1",
			CompMain:        7200,
			CompPlus:        14400,
			CompAll:         21600,
			Similarity:      0.5,
		},
		{
			GameID:          2,
			GameName:        "Game2",
			ProfilePlatform: "Platform2",
			GameImage:       "Image2",
			CompMain:        3600,
			CompPlus:        7200,
			CompAll:         10800,
			Similarity:      0.7,
		},
	}

	searchGame := &SearchGame{Data: data}
	expected := []*SearchGameSimple{
		{
			GameID:          1,
			GameName:        "Game1",
			ProfilePlatform: "Platform1",
			GameImage:       "Image1",
			CompMain:        math.Round(float64(7200) / 3600),
			CompPlus:        math.Round(float64(14400) / 3600),
			CompAll:         math.Round(float64(21600) / 3600),
			Similarity:      0.5,
		},
		{
			GameID:          2,
			GameName:        "Game2",
			ProfilePlatform: "Platform2",
			GameImage:       "Image2",
			CompMain:        math.Round(float64(3600) / 3600),
			CompPlus:        math.Round(float64(7200) / 3600),
			CompAll:         math.Round(float64(10800) / 3600),
			Similarity:      0.7,
		},
	}

	got := searchGame.Reduce()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Reduce() = %v, want %v", got, expected)
	}
}

func Test_Reduce(t *testing.T) {
	s := &GameDetails{
		Props: struct {
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
		}{
			PageProps: struct {
				Game struct {
					Count int
					Data  GameDetailsGameData
				}
				IgnWikiSlug  string
				IgnMap       *GameDetailsIgnMap
				IgnWikiNav   []GameDetailsIgnWikiNav
				PageMetadata GameDetailsPageMetadata
			}{
				Game: struct {
					Count int
					Data  GameDetailsGameData
				}{
					Count: 1,
					Data: GameDetailsGameData{
						Game: []GameDetailsGameDataGame{
							{
								GameID:          1,
								GameName:        "test game",
								ProfilePlatform: "test platform",
								GameImage:       "test image",
								CompMain:        3600,
								CompPlus:        3600,
								CompAll:         3600,
							},
						},
					},
				},
			},
		},
	}

	expectedResult := &GameDetailSimple{
		GameID:          1,
		GameName:        "test game",
		ProfilePlatform: "test platform",
		GameImage:       "test image",
		CompMain:        math.Round(float64(3600) / 3600),
		CompPlus:        math.Round(float64(3600) / 3600),
		CompAll:         math.Round(float64(3600) / 3600),
	}

	reduced := s.Reduce()

	if *reduced != *expectedResult {
		t.Errorf("Reduce() = %v, want %v", *reduced, *expectedResult)
	}
}

func Test_SearchSimple(t *testing.T) {
	ctx := t.Context()

	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	result, err := mockClient.SearchSimple(ctx, "The Witcher 3: Wild Hunt", SearchModifierNone)
	if err != nil {
		t.Fatalf("SearchSimple() = %v", err)
	}

	if result[0].GameID != 10270 {
		t.Errorf("Search() gameID = %v, want %v", result[0].GameID, 10270)
	}
}

func Test_DetailsSimple(t *testing.T) {
	ctx := t.Context()

	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	result, err := mockClient.DetailSimple(ctx, 10270)
	if err != nil {
		t.Fatalf("DetailSimple() = %v", err)
	}

	if result.GameID != 10270 {
		t.Errorf("DetailSimple() gameID = %v, want %v", result.GameID, 10270)
	}
}
