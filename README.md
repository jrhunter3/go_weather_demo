# go_weather_demo

A CLI weather tool written in Go that fetches real-time weather data from the [Open-Meteo](https://open-meteo.com) public API. Demonstrates idiomatic Go patterns including HTTP clients, JSON parsing, CLI flag handling, concurrency (goroutines), error handling, and table-formatted output.

## Prerequisites

- **Go 1.22+** — [Download](https://go.dev/dl/)

## Setup

```bash
git clone <repo-url>
cd go_weather_demo
go build ./...
```

## Usage

```bash
# Fetch weather for one or more cities
go run . "New York" "Los Angeles" "Chicago"

# Use metric units
go run . -units metric "London"

# Set a custom request timeout
go run . -timeout 5s "Tokyo"
```

### Flags

| Flag      | Default     | Description                |
|-----------|-------------|----------------------------|
| `-units`  | `imperial`  | `imperial` or `metric`     |
| `-timeout`| `10s`       | Request timeout duration   |

## Testing

```bash
go test ./... -v
```

## Project structure

```
.
├── main.go              # CLI entry point
├── weather/
│   ├── weather.go       # Fetch & parse weather data
│   └── weather_test.go  # Unit tests with httptest
├── output/
│   ├── output.go        # Formatting & display
│   └── output_test.go
├── go.mod
├── AGENTS.md
└── PLAN.md
```
