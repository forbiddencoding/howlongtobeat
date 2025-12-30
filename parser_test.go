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

func Test_htmlScriptDataParserByID(t *testing.T) {
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
	parseFunc := mockClient.htmlScriptDataParserByID(&dest, "__NEXT_DATA__")

	if err = parseFunc(resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dest.Props.PageProps.Game.Data.Game[0].GameID != 10270 {
		t.Fatalf("unexpected game id: %v", dest.Props.PageProps.Game.Data.Game[0].GameID)
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

func Benchmark_htmlScriptDataParserByID(b *testing.B) {
	htmlFile, err := os.Open("test_files/test_html_parser.html")
	if err != nil {
		b.Fatalf("error opening HTML test file: %v", err)
	}
	defer htmlFile.Close()

	mockData, err := io.ReadAll(htmlFile)
	if err != nil {
		b.Fatalf("error reading HTML test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}

	dest := gameDetailsResponse{}
	mockClient := &Client{}
	parseFunc := mockClient.htmlScriptDataParserByID(&dest, "__NEXT_DATA__")

	if err = parseFunc(resp); err != nil {
		b.Fatalf("unexpected error: %v", err)
	}
}

func Benchmark_nextDataParser(b *testing.B) {
	htmlFile, err := os.Open("test_files/test_html_parser.html")
	if err != nil {
		b.Fatalf("error opening HTML test file: %v", err)
	}
	defer htmlFile.Close()

	mockData, err := io.ReadAll(htmlFile)
	if err != nil {
		b.Fatalf("error reading HTML test file: %v", err)
	}

	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(mockData)
	}))
	defer rs.Close()

	resp, err := http.Get(rs.URL)
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}

	dest := gameDetailsResponse{}
	mockClient := &Client{}
	parseFunc := mockClient.nextDataParser(&dest)

	if err = parseFunc(resp); err != nil {
		b.Fatalf("unexpected error: %v", err)
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
