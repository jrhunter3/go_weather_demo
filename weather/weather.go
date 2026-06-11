package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Conditions struct {
	City      string
	TempC     float64
	TempF     float64
	Condition string
	Humidity  float64
	WindKph   float64
	WindMph   float64
	FetchedAt time.Time
}

type Result struct {
	Conditions Conditions
	Err        error
}

type Client struct {
	httpClient  *http.Client
	weatherBase string
	geocodeBase string
}

func NewClient() *Client {
	return &Client{
		httpClient:  http.DefaultClient,
		weatherBase: "https://api.open-meteo.com/v1",
		geocodeBase: "https://geocoding-api.open-meteo.com/v1",
	}
}

type geoResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type geoResponse struct {
	Results []geoResult `json:"results"`
}

type currentWeather struct {
	Temperature  float64 `json:"temperature_2m"`
	Humidity     float64 `json:"relative_humidity_2m"`
	WeatherCode  int     `json:"weather_code"`
	WindSpeed    float64 `json:"wind_speed_10m"`
}

type weatherResponse struct {
	Current currentWeather `json:"current"`
}

func (c *Client) geocode(ctx context.Context, city string) (float64, float64, error) {
	u := fmt.Sprintf("%s/search?name=%s&count=1&language=en&format=json", c.geocodeBase, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("creating geocode request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("geocode request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocode API returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("reading geocode response: %w", err)
	}
	var geo geoResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		return 0, 0, fmt.Errorf("decoding geocode response: %w", err)
	}
	if len(geo.Results) == 0 {
		return 0, 0, fmt.Errorf("city %q not found", city)
	}
	return geo.Results[0].Latitude, geo.Results[0].Longitude, nil
}

func (c *Client) fetchCurrent(ctx context.Context, lat, lon float64) (currentWeather, error) {
	u := fmt.Sprintf("%s/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m", c.weatherBase, lat, lon)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return currentWeather{}, fmt.Errorf("creating weather request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return currentWeather{}, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return currentWeather{}, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return currentWeather{}, fmt.Errorf("reading weather response: %w", err)
	}
	var w weatherResponse
	if err := json.Unmarshal(body, &w); err != nil {
		return currentWeather{}, fmt.Errorf("decoding weather response: %w", err)
	}
	return w.Current, nil
}

var weatherCodes = map[int]string{
	0:    "Clear",
	1:    "Mainly clear",
	2:    "Partly cloudy",
	3:    "Overcast",
	45:   "Foggy",
	48:   "Depositing rime fog",
	51:   "Light drizzle",
	53:   "Moderate drizzle",
	55:   "Dense drizzle",
	56:   "Light freezing drizzle",
	57:   "Dense freezing drizzle",
	61:   "Slight rain",
	63:   "Moderate rain",
	65:   "Heavy rain",
	66:   "Light freezing rain",
	67:   "Heavy freezing rain",
	71:   "Slight snow",
	73:   "Moderate snow",
	75:   "Heavy snow",
	77:   "Snow grains",
	80:   "Slight rain showers",
	81:   "Moderate rain showers",
	82:   "Violent rain showers",
	85:   "Slight snow showers",
	86:   "Heavy snow showers",
	95:   "Thunderstorm",
	96:   "Thunderstorm with slight hail",
	99:   "Thunderstorm with heavy hail",
}

func conditionText(code int) string {
	if s, ok := weatherCodes[code]; ok {
		return s
	}
	return "Unknown"
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

func (c *Client) Get(ctx context.Context, city string) (Conditions, error) {
	lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return Conditions{}, err
	}
	w, err := c.fetchCurrent(ctx, lat, lon)
	if err != nil {
		return Conditions{}, err
	}
	return Conditions{
		City:      city,
		TempC:     round1(w.Temperature),
		TempF:     round1(w.Temperature*9/5 + 32),
		Condition: conditionText(w.WeatherCode),
		Humidity:  round1(w.Humidity),
		WindKph:   round1(w.WindSpeed),
		WindMph:   round1(w.WindSpeed * 0.621371),
		FetchedAt: time.Now(),
	}, nil
}

func (c *Client) GetMany(ctx context.Context, cities []string) []Result {
	results := make([]Result, len(cities))
	var wg sync.WaitGroup
	for i, city := range cities {
		wg.Add(1)
		i, city := i, city
		go func() {
			defer wg.Done()
			cond, err := c.Get(ctx, city)
			results[i] = Result{Conditions: cond, Err: err}
		}()
	}
	wg.Wait()
	return results
}
