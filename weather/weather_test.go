package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestGetForecast_Success(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	fr, err := client.GetForecast(context.Background(), "New York", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fr.City != "New York" {
		t.Errorf("expected City 'New York', got %q", fr.City)
	}
	if len(fr.Days) != 3 {
		t.Fatalf("expected 3 forecast days, got %d", len(fr.Days))
	}

	if fr.Days[0].Date != "2025-01-20" {
		t.Errorf("expected date '2025-01-20', got %q", fr.Days[0].Date)
	}
	if fr.Days[0].TempMaxC != 10.0 {
		t.Errorf("expected TempMaxC 10.0, got %.1f", fr.Days[0].TempMaxC)
	}
	if fr.Days[0].TempMaxF != 50.0 {
		t.Errorf("expected TempMaxF 50.0, got %.1f", fr.Days[0].TempMaxF)
	}
	if fr.Days[0].Condition != "Slight rain" {
		t.Errorf("expected Condition 'Slight rain', got %q", fr.Days[0].Condition)
	}
}

func TestGetForecast_UnknownCity(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	_, err := client.GetForecast(context.Background(), "NonexistentCityXYZ", 3)
	if err == nil {
		t.Fatal("expected error for unknown city, got nil")
	}
}

func TestGetManyForecasts(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	results := client.GetManyForecasts(context.Background(), []string{"New York", "Los Angeles", "NonexistentCityXYZ"}, 3)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if results[0].Err != nil {
		t.Errorf("unexpected error for New York: %v", results[0].Err)
	}
	if len(results[0].Days) != 3 {
		t.Errorf("expected 3 forecast days for New York, got %d", len(results[0].Days))
	}
	if results[2].Err == nil {
		t.Error("expected error for unknown city, got nil")
	}
}

func TestGet_ContextCancel(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Get(ctx, "New York")
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestGetMany_ContextCancel(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results := client.GetMany(ctx, []string{"New York", "Los Angeles"})
	for i, r := range results {
		if r.Err == nil {
			t.Errorf("result[%d]: expected error for cancelled context, got nil", i)
		}
	}
}

func TestGetForecast_ContextCancel(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	client := &Client{
		httpClient:  srv.Client(),
		weatherBase: srv.URL + "/v1",
		geocodeBase: srv.URL + "/v1",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetForecast(ctx, "New York", 3)
	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("daily") {
			resp := map[string]any{
				"daily": map[string]any{
					"time":               []string{"2025-01-20", "2025-01-21", "2025-01-22"},
					"temperature_2m_max": []float64{10.0, 12.5, 8.0},
					"temperature_2m_min": []float64{2.0, 5.0, 1.0},
					"weather_code":       []int{61, 3, 80},
					"precipitation_sum":  []float64{2.5, 0.0, 1.2},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}
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
					"temperature_2m":       15.0,
					"relative_humidity_2m": 72.0,
					"weather_code":         0,
					"wind_speed_10m":       12.5,
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestConditionText(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, "Clear"},
		{1, "Mainly clear"},
		{61, "Slight rain"},
		{99, "Thunderstorm with heavy hail"},
		{999, "Unknown"},
	}
	for _, tc := range tests {
		got := conditionText(tc.code)
		if got != tc.want {
			t.Errorf("conditionText(%d) = %q, want %q", tc.code, got, tc.want)
		}
	}
}

func TestRound1(t *testing.T) {
	tests := []struct {
		v    float64
		want float64
	}{
		{15.0, 15.0},
		{12.345, 12.3},
		{12.3456, 12.3},
		{0.621371, 0.6},
	}
	for _, tc := range tests {
		got := round1(tc.v)
		if got != tc.want {
			t.Errorf("round1(%f) = %f, want %f", tc.v, got, tc.want)
		}
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.httpClient != http.DefaultClient {
		t.Error("expected http.DefaultClient")
	}
	if !strings.HasPrefix(c.weatherBase, "https://") {
		t.Errorf("expected HTTPS weatherBase, got %q", c.weatherBase)
	}
	if !strings.HasPrefix(c.geocodeBase, "https://") {
		t.Errorf("expected HTTPS geocodeBase, got %q", c.geocodeBase)
	}
}
