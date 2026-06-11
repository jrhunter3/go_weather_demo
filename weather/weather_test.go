package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet_Success(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	ctx := context.Background()
	cond, err := client.Get(ctx, "New York")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cond.City != "New York" {
		t.Errorf("expected City 'New York', got %q", cond.City)
	}
	if cond.TempC != 15.0 {
		t.Errorf("expected TempC 15.0, got %.1f", cond.TempC)
	}
	if cond.TempF != 59.0 {
		t.Errorf("expected TempF 59.0, got %.1f", cond.TempF)
	}
	if cond.Condition != "Clear" {
		t.Errorf("expected Condition 'Clear', got %q", cond.Condition)
	}
	if cond.Humidity != 72.0 {
		t.Errorf("expected Humidity 72.0, got %.1f", cond.Humidity)
	}
	if cond.WindKph != 12.5 {
		t.Errorf("expected WindKph 12.5, got %.1f", cond.WindKph)
	}
	if cond.FetchedAt.IsZero() {
		t.Error("expected FetchedAt to be set")
	}
}

func TestGet_UnknownCity(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	_, err := client.Get(context.Background(), "NonexistentCityXYZ")
	if err == nil {
		t.Fatal("expected error for unknown city, got nil")
	}
}

func TestGetMany(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	results := client.GetMany(context.Background(), []string{"New York", "Los Angeles", "NonexistentCityXYZ"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0].Err != nil {
		t.Errorf("unexpected error for New York: %v", results[0].Err)
	}
	if results[1].Err != nil {
		t.Errorf("unexpected error for Los Angeles: %v", results[1].Err)
	}
	if results[2].Err == nil {
		t.Error("expected error for unknown city, got nil")
	}
}

func TestGet_APIDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	_, err := client.Get(context.Background(), "New York")
	if err == nil {
		t.Fatal("expected error for API down, got nil")
	}
}

func TestGet_EmptyCity(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	_, err := client.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty city, got nil")
	}
}

func TestGetMany_EmptyInput(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	results := client.GetMany(context.Background(), nil)
	if len(results) != 0 {
		t.Errorf("expected 0 results for nil input, got %d", len(results))
	}
}

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/search":
			city := r.URL.Query().Get("name")
			if city == "" || city == "NonexistentCityXYZ" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"results":[]}`)
				return
			}
			resp := map[string]any{
				"results": []map[string]any{
					{
						"name":      city,
						"latitude":  40.71,
						"longitude": -74.01,
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)

		case "/v1/forecast":
			resp := map[string]any{
				"current": map[string]any{
					"temperature_2m":      15.0,
					"relative_humidity_2m": 72.0,
					"weather_code":        0,
					"wind_speed_10m":      12.5,
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
