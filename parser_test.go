package howlongtobeat

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_jsonParser(t *testing.T) {
	jsonFile, err := os.Open("test_files/test_json_parser.json")
	if err != nil {
		t.Fatalf("error opening JSON test file: %v", err)
	}
	defer jsonFile.Close()

	mockData, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("error reading JSON test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockClient := &Client{}
	dest := SearchGame{}
	parseFunc := mockClient.jsonParser(&dest)

	if err = parseFunc(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dest.Data[0].GameName != "The Witcher 3: Wild Hunt" {
		t.Fatalf("unexpected game name: %v", dest.Data[0].GameName)
	}

	if dest.Data[0].GameID != 10270 {
		t.Fatalf("unexpected game id: %v", dest.Data[0].GameID)
	}

}

func Test_nextDataParser(t *testing.T) {
	htmlFile, err := os.Open("test_files/test_html_parser.html")
	if err != nil {
		t.Fatalf("error opening HTML test file: %v", err)
	}
	defer htmlFile.Close()

	mockData, err := io.ReadAll(htmlFile)
	if err != nil {
		t.Fatalf("error reading HTML test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dest := gameDetailsResponse{}
	mockClient := &Client{}
	parseFunc := mockClient.nextDataParser(&dest)

	if err = parseFunc(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dest.Props.PageProps.Game.Data.Game[0].GameID != 10270 {
		t.Fatalf("unexpected game id: %v", dest.Props.PageProps.Game.Data.Game[0].GameID)
	}
}

func Test_scriptPathParser(t *testing.T) {
	htmlFile, err := os.Open("test_files/test_html_parser.html")
	if err != nil {
		t.Fatalf("error opening HTML test file: %v", err)
	}
	defer htmlFile.Close()

	mockData, err := io.ReadAll(htmlFile)
	if err != nil {
		t.Fatalf("error reading HTML test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apiData := ApiData{}
	mockClient := &Client{}
	parseFunc := mockClient.scriptParser(&apiData)

	if err = parseFunc(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(apiData.scriptPaths) != 14 {
		t.Fatalf("unexpected script path: %v", apiData.scriptPaths)
	}
}

func Test_endpointParser(t *testing.T) {
	jsFile, err := os.Open("test_files/test_endpoint_parser.js")
	if err != nil {
		t.Fatalf("error opening JS test file: %v", err)
	}
	defer jsFile.Close()

	mockData, err := io.ReadAll(jsFile)
	if err != nil {
		t.Fatalf("error reading JS test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	apiData := ApiData{}
	mockClient := &Client{}
	parseFunc := mockClient.endpointParser(&apiData)

	if err = parseFunc(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := "/api/search"
	if apiData.endpointPath != expectedPath {
		t.Fatalf("unexpected endpoint path: %s, expected: %s", apiData.endpointPath, expectedPath)
	}
}
