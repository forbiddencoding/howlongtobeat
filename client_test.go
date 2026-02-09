package howlongtobeat

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}
}

func ExampleNew() {
	_, err := New()
	if err != nil {
		panic(err)
	}
}

func TestWithDefaultRequestTimeout(t *testing.T) {
	mockClient, err := New()
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if mockClient.client.Timeout != defaultRequestTimeout {
		t.Fatalf("WithRequestTimeout() did not set the default timeout")
	}
}

func TestWithCustomZeroRequestTimeout(t *testing.T) {
	mockClient, err := New(WithRequestTimeout(0))
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if mockClient.client.Timeout != defaultRequestTimeout {
		t.Fatalf("WithRequestTimeout() did not set the default timeout")
	}
}

func TestWithCustomNonZeroRequestTimeout(t *testing.T) {
	mockClient, err := New(WithRequestTimeout(10))
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if mockClient.client.Timeout != 10*time.Second {
		t.Fatalf("WithRequestTimeout() did not set the custom timeout")
	}
}

func TestWithCustomHTTPClient(t *testing.T) {
	customHTTPClient := http.DefaultClient
	mockClient, err := New(WithHTTPClient(customHTTPClient))
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if mockClient.client != customHTTPClient {
		t.Fatalf("WithHTTPClient() did not set the custom HTTP client")
	}
}

func TestNewWithOptions(t *testing.T) {
	customHTTPClient := http.DefaultClient

	mockClient, err := New(
		WithHTTPClient(customHTTPClient),
		WithRequestTimeout(10),
	)
	if err != nil {
		t.Fatalf("New() returned error: %v", err)
	}

	if mockClient.client.Timeout != 10*time.Second {
		t.Fatalf("WithRequestTimeout() did not set the custom timeout")
	}

	if mockClient.client != customHTTPClient {
		t.Fatalf("WithHTTPClient() did not set the custom HTTP client")
	}
}

func Test_do(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	mockClient := server.Client()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	responseParser := func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		if !strings.Contains(string(body), "OK") {
			return errors.New("unexpected body content")
		}

		return nil
	}

	c := Client{
		client: mockClient,
	}

	err = c.do(req, responseParser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_do_Invalid_Status(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("500"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	mockClient := server.Client()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	responseParser := func(r *http.Response) error {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		if !strings.Contains(string(body), "OK") {
			return errors.New("unexpected body content")
		}

		return nil
	}

	c := Client{
		client: mockClient,
	}

	err = c.do(req, responseParser)
	if strings.Compare(err.Error(), "unexpected status code: 500") != 0 {
		t.Fatalf("unexpected error: %v", err)
	}
}
