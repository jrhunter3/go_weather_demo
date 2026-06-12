package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"go_weather_demo/weather"
)

func TestPrint_Imperial(t *testing.T) {
	var b strings.Builder
	results := []weather.Result{
		{Conditions: weather.Conditions{City: "New York", TempF: 59, TempC: 15, Condition: "Clear", Humidity: 72, WindMph: 7.8, WindKph: 12.5}},
		{Conditions: weather.Conditions{City: "Los Angeles", TempF: 75, TempC: 24, Condition: "Mainly clear", Humidity: 45, WindMph: 8.1, WindKph: 13.0}},
	}

	useColor = false
	Print(&b, results, Imperial, 450*time.Millisecond)
	got := b.String()

	if !strings.Contains(got, "New York") {
		t.Errorf("expected output to contain 'New York', got:\n%s", got)
	}
	if !strings.Contains(got, "59°F") {
		t.Errorf("expected output to contain '59°F', got:\n%s", got)
	}
	if !strings.Contains(got, "Clear") {
		t.Errorf("expected output to contain 'Clear', got:\n%s", got)
	}
	if !strings.Contains(got, "72%") {
		t.Errorf("expected output to contain '72%%', got:\n%s", got)
	}
	if !strings.Contains(got, "7.8 mph") {
		t.Errorf("expected output to contain '7.8 mph', got:\n%s", got)
	}
	if !strings.Contains(got, "Fetched 2 cities in 0.45s") {
		t.Errorf("expected summary line 'Fetched 2 cities in 0.45s', got:\n%s", got)
	}
}

func TestPrint_Metric(t *testing.T) {
	var b strings.Builder
	results := []weather.Result{
		{Conditions: weather.Conditions{City: "Osaka", TempF: 77, TempC: 25, Condition: "Partly cloudy", Humidity: 60, WindMph: 6.2, WindKph: 10.0}},
	}

	useColor = false
	Print(&b, results, Metric, 300*time.Millisecond)
	got := b.String()

	if !strings.Contains(got, "25°C") {
		t.Errorf("expected output to contain '25°C', got:\n%s", got)
	}
	if !strings.Contains(got, "10.0 km/h") {
		t.Errorf("expected output to contain '10.0 km/h', got:\n%s", got)
	}
	if !strings.Contains(got, "Fetched 1 city in 0.30s") {
		t.Errorf("expected singular summary, got:\n%s", got)
	}
}

func TestPrint_ErrorRow(t *testing.T) {
	var b strings.Builder
	results := []weather.Result{
		{Conditions: weather.Conditions{City: "Atlantis"}, Err: fmt.Errorf("city not found")},
	}

	useColor = false
	Print(&b, results, Imperial, 100*time.Millisecond)
	got := b.String()

	if !strings.Contains(got, "Atlantis") {
		t.Errorf("expected output to contain 'Atlantis', got:\n%s", got)
	}
	if !strings.Contains(got, "not found") {
		t.Errorf("expected error message in output, got:\n%s", got)
	}
}

func TestPrint_Empty(t *testing.T) {
	var b strings.Builder
	useColor = false
	Print(&b, nil, Imperial, 0)
	got := b.String()

	if !strings.Contains(got, "Fetched 0 cities in 0.00s") {
		t.Errorf("expected 'Fetched 0 cities in 0.00s', got:\n%s", got)
	}
}

func TestPrintForecast_Success(t *testing.T) {
	var b strings.Builder
	results := []weather.ForecastResult{
		{
			City: "New York",
			Days: []weather.DailyForecast{
				{Date: "2025-01-20", TempMaxC: 10.0, TempMaxF: 50.0, TempMinC: 2.0, TempMinF: 35.6, Condition: "Slight rain", PrecipMM: 2.5},
				{Date: "2025-01-21", TempMaxC: 12.5, TempMaxF: 54.5, TempMinC: 5.0, TempMinF: 41.0, Condition: "Overcast", PrecipMM: 0.0},
			},
		},
	}

	useColor = false
	PrintForecast(&b, results, Imperial, 600*time.Millisecond)
	got := b.String()

	if !strings.Contains(got, "New York") {
		t.Errorf("expected output to contain 'New York', got:\n%s", got)
	}
	if !strings.Contains(got, "50°F") {
		t.Errorf("expected output to contain '50°F', got:\n%s", got)
	}
	if !strings.Contains(got, "Slight rain") {
		t.Errorf("expected output to contain 'Slight rain', got:\n%s", got)
	}
	if !strings.Contains(got, "2.5 mm") {
		t.Errorf("expected output to contain '2.5 mm', got:\n%s", got)
	}
	if !strings.Contains(got, "Fetched 1 city in 0.60s") {
		t.Errorf("expected summary line, got:\n%s", got)
	}
}

func TestPrintForecast_MultipleCities(t *testing.T) {
	var b strings.Builder
	results := []weather.ForecastResult{
		{City: "New York", Days: []weather.DailyForecast{{Date: "2025-01-20", TempMaxC: 10.0, TempMaxF: 50.0, TempMinC: 2.0, TempMinF: 35.6, Condition: "Rain", PrecipMM: 2.5}}},
		{City: "Atlantis", Err: fmt.Errorf("city not found")},
	}

	useColor = false
	PrintForecast(&b, results, Metric, 500*time.Millisecond)
	got := b.String()

	if !strings.Contains(got, "New York") {
		t.Errorf("expected 'New York', got:\n%s", got)
	}
	if !strings.Contains(got, "Atlantis") {
		t.Errorf("expected 'Atlantis', got:\n%s", got)
	}
	if !strings.Contains(got, "not found") {
		t.Errorf("expected error message, got:\n%s", got)
	}
	if !strings.Contains(got, "Fetched 2 cities in 0.50s") {
		t.Errorf("expected summary, got:\n%s", got)
	}
}

func TestPrintJSON_Results(t *testing.T) {
	var b strings.Builder
	results := []weather.Result{
		{Conditions: weather.Conditions{City: "New York", TempF: 59, TempC: 15, Condition: "Clear", Humidity: 72, WindMph: 7.8, WindKph: 12.5}},
	}

	err := PrintJSON(&b, results)
	if err != nil {
		t.Fatalf("PrintJSON error: %v", err)
	}

	var decoded []weather.Result
	if err := json.Unmarshal([]byte(b.String()), &decoded); err != nil {
		t.Fatalf("JSON unmarshal error: %v\nbody: %s", err, b.String())
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 result, got %d", len(decoded))
	}
	if decoded[0].Conditions.City != "New York" {
		t.Errorf("expected City 'New York', got %q", decoded[0].Conditions.City)
	}
}

func TestPrintJSON_Forecast(t *testing.T) {
	var b strings.Builder
	results := []weather.ForecastResult{
		{City: "Tokyo", Days: []weather.DailyForecast{{Date: "2025-01-20", TempMaxC: 15.0, TempMaxF: 59.0, TempMinC: 8.0, TempMinF: 46.4, Condition: "Clear", PrecipMM: 0.0}}},
	}

	err := PrintJSON(&b, results)
	if err != nil {
		t.Fatalf("PrintJSON error: %v", err)
	}

	var decoded []weather.ForecastResult
	if err := json.Unmarshal([]byte(b.String()), &decoded); err != nil {
		t.Fatalf("JSON unmarshal error: %v\nbody: %s", err, b.String())
	}
	if len(decoded) != 1 {
		t.Fatalf("expected 1 result, got %d", len(decoded))
	}
	if decoded[0].City != "Tokyo" {
		t.Errorf("expected City 'Tokyo', got %q", decoded[0].City)
	}
	if len(decoded[0].Days) != 1 {
		t.Fatalf("expected 1 day, got %d", len(decoded[0].Days))
	}
}
