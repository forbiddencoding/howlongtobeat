package howlongtobeat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func Test_detailHTTPRequest(t *testing.T) {
	ctx := t.Context()

	var (
		headers = map[string]string{
			http.CanonicalHeaderKey("User-Agent"):       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
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
		gameID     = 14286
		mockClient = &Client{}
	)

	req, err := mockClient.detailHTTPRequest(ctx, gameID)
	if err != nil {
		t.Fatalf("received error from detailHTTPRequest(): %v", err)
	}

	if req.Method != http.MethodGet {
		t.Fatalf("detailHTTPRequest() did not set the correct method: want: %s, received: %s", http.MethodGet, req.Method)
	}

	if req.URL.String() != fmt.Sprintf("%s/%d", hltbGameURL, gameID) {
		t.Fatalf("detailHTTPRequest() did not set the correct URL: want: %s, received: %s", fmt.Sprintf("%s/%d", hltbGameURL, gameID), req.URL.String())
	}

	if req.Body != nil {
		t.Fatalf("detailHTTPRequest() body should be nil, received: %v", req.Body)
	}

	var errs error

	for header, value := range headers {
		if req.Header.Get(header) != value {
			errs = errors.Join(errs, fmt.Errorf("detailHTTPRequest() did not set the correct %s header: want: %s, received: %s", header, value, req.Header.Get(header)))
		}
	}

	if errs != nil {
		t.Fatal(errs)
	}
}

func Test_Detail(t *testing.T) {
	ctx := t.Context()

	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	result, err := mockClient.Detail(ctx, 10270)
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}

	if result.Props.PageProps.Game.Data.Game[0].GameID != 10270 {
		t.Fatalf("Detail() did not return the correct game: want: %d, received: %d", 10270, result.Props.PageProps.Game.Data.Game[0].GameID)
	}
}

func Test_Detail_EmptyIgnWikiNav(t *testing.T) {
	ctx := t.Context()

	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	result, err := mockClient.Detail(ctx, 14286)
	if err != nil {
		t.Fatalf("Detail() error = %v", err)
	}

	if result.Props.PageProps.Game.Data.Game[0].GameID != 14286 {
		t.Fatalf("Detail() did not return the correct game: want: %d, received: %d", 14286, result.Props.PageProps.Game.Data.Game[0].GameID)
	}
}

func Test_Detail_InvalidGameID(t *testing.T) {
	ctx := t.Context()
	mockClient := &Client{}

	_, err := mockClient.Detail(ctx, 0)
	if err == nil {
		t.Fatalf("Detail() expected error, but received: %v", err)
	}
	if !errors.Is(err, GameIDRequiredErr) {
		t.Fatalf(`Detail() expected "%v" error, but received: %v`, GameIDRequiredErr, err)
	}
}

func Test_convertResponseToGameDetails_IgnNav(t *testing.T) {
	// Construct a gameDetailsResponse object
	var dummyIgnNav = map[string]string{"__typename": "WikiNavigation", "label": "The Witcher 3 Guide", "url": ""}
	dummyIgnNavJSON, _ := json.Marshal(dummyIgnNav)

	g := &gameDetailsResponse{}
	g.Props.PageProps.IgnWikiNav = dummyIgnNavJSON

	result, err := g.convertResponseToGameDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if result.Props.PageProps.IgnWikiNav == nil {
		t.Error("Expected IgnNav to not be nil")
	}

	if result.Props.PageProps.IgnWikiNav[0].Label != "The Witcher 3 Guide" {
		t.Error("IgnMap not correctly converted")
	}
}

func Test_convertResponseToGameDetails_IgnMap(t *testing.T) {
	var dummyIgnMap = &GameDetailsIgnMap{
		Typename:    "Map",
		MapType:     "mapgenie",
		MapName:     "White Orchard",
		MapSlug:     "white-orchard",
		ObjectSlug:  "the-witcher-3-wild-hunt",
		ObjectName:  "The Witcher 3",
		Width:       1,
		Height:      1,
		MinZoom:     8,
		MaxZoom:     14,
		Tilesets:    []string{"https://tiles.mapgenie.io/games/witcher-3/white-orchard/default/{z}/{x}/{y}.png"},
		InitialLat:  83.937238401332,
		InitialLng:  -168.44211701243,
		InitialZoom: 11,
	}

	var dummyIgnNav = map[string]string{"__typename": "WikiNavigation", "label": "The Witcher 3 Guide", "url": ""}
	dummyIgnNavJSON, _ := json.Marshal(dummyIgnNav)

	g := &gameDetailsResponse{}
	g.Props.PageProps.IgnWikiNav = dummyIgnNavJSON
	g.Props.PageProps.IgnMap = dummyIgnMap

	result, err := g.convertResponseToGameDetails()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if result.Props.PageProps.IgnMap == nil {
		t.Error("Expected IgnMap to not be nil")
	}

	if result.Props.PageProps.IgnMap.MapName != "White Orchard" {
		t.Error("IgnMap not correctly converted")
	}
}
