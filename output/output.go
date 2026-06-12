package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"go_weather_demo/weather"
)

// useColor controls ANSI color output. Disabled when NO_COLOR is set.
var useColor = os.Getenv("NO_COLOR") == ""

// Units represents the temperature unit system for display.
type Units int

// Unit system constants.
const (
	Imperial Units = iota
	Metric
)

// Print writes a formatted table of current weather results to w.
func Print(w io.Writer, results []weather.Result, units Units, elapsed time.Duration) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	fmt.Fprintln(tw, "  City\tTemp\tCondition\tHumidity\tWind")
	fmt.Fprintln(tw, "  ─────────────────────────────────────────────")

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(tw, "  %s\t—\t%s\t—\t—\n", r.Conditions.City, colorError(r.Err.Error()))
			continue
		}
		c := r.Conditions
		temp := formatTemp(c, units)
		wind := formatWind(c, units)
		cond := colorCondition(c.Condition, c.Condition)
		humi := fmt.Sprintf("%.0f%%", c.Humidity)
		fmt.Fprintf(tw, "  %s\t%s\t%s\t%s\t%s\n", c.City, temp, cond, humi, wind)
	}

	tw.Flush()

	n := len(results)
	label := "city"
	if n != 1 {
		label = "cities"
	}
	summary := fmt.Sprintf("\n  Fetched %d %s in %.2fs\n", n, label, elapsed.Seconds())
	fmt.Fprint(w, summary)
}

// PrintForecast writes a formatted forecast table to w.
func PrintForecast(w io.Writer, results []weather.ForecastResult, units Units, elapsed time.Duration) {
	for i, r := range results {
		if i > 0 {
			fmt.Fprintln(w)
		}
		if r.Err != nil {
			fmt.Fprintf(w, "  %s: %s\n", r.City, colorError(r.Err.Error()))
			continue
		}
		fmt.Fprintf(w, "  %s:\n", r.City)
		tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
		fmt.Fprintln(tw, "  Date\tMax\tMin\tCondition\tPrecip")
		fmt.Fprintln(tw, "  ─────────────────────────────────────────")
		for _, d := range r.Days {
			maxStr := colorTemp(formatTempVal(d.TempMaxC, d.TempMaxF, units), maxTemp(d, units))
			minStr := formatTempVal(d.TempMinC, d.TempMinF, units)
			cond := colorCondition(d.Condition, d.Condition)
			precip := fmt.Sprintf("%.1f mm", d.PrecipMM)
			fmt.Fprintf(tw, "  %s\t%s\t%s\t%s\t%s\n", d.Date, maxStr, minStr, cond, precip)
		}
		tw.Flush()
	}

	n := len(results)
	label := "city"
	if n != 1 {
		label = "cities"
	}
	summary := fmt.Sprintf("\n  Fetched %d %s in %.2fs\n", n, label, elapsed.Seconds())
	fmt.Fprint(w, summary)
}

// PrintJSON writes v as indented JSON to w.
func PrintJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func formatTemp(c weather.Conditions, units Units) string {
	switch units {
	case Metric:
		return colorTemp(fmt.Sprintf("%.0f°C", c.TempC), c.TempC)
	default:
		return colorTemp(fmt.Sprintf("%.0f°F", c.TempF), c.TempF)
	}
}

func formatTempVal(tempC, tempF float64, units Units) string {
	switch units {
	case Metric:
		return fmt.Sprintf("%.0f°C", tempC)
	default:
		return fmt.Sprintf("%.0f°F", tempF)
	}
}

func maxTemp(d weather.DailyForecast, units Units) float64 {
	switch units {
	case Metric:
		return d.TempMaxC
	default:
		return d.TempMaxF
	}
}

func formatWind(c weather.Conditions, units Units) string {
	switch units {
	case Metric:
		return fmt.Sprintf("%.1f km/h", c.WindKph)
	default:
		return fmt.Sprintf("%.1f mph", c.WindMph)
	}
}

func colorTemp(s string, t float64) string {
	if !useColor {
		return s
	}
	switch {
	case t >= 100:
		return red(s)
	case t >= 80:
		return yellow(s)
	case t <= 32:
		return cyan(s)
	default:
		return s
	}
}

func colorCondition(cond, s string) string {
	if !useColor {
		return s
	}
	switch {
	case strings.Contains(cond, "Clear"), strings.Contains(cond, "Partly"):
		return yellow(s)
	case strings.Contains(cond, "Cloudy"), strings.Contains(cond, "Overcast"), strings.Contains(cond, "Fog"):
		return dim(s)
	case strings.Contains(cond, "Rain"), strings.Contains(cond, "Drizzle"):
		return blue(s)
	case strings.Contains(cond, "Snow"):
		return white(s)
	case strings.Contains(cond, "Thunderstorm"):
		return red(s)
	default:
		return s
	}
}

func colorError(s string) string {
	if !useColor {
		return s
	}
	return red(s)
}

func red(s string) string    { return "\033[31m" + s + "\033[0m" }
func yellow(s string) string { return "\033[33m" + s + "\033[0m" }
func blue(s string) string   { return "\033[34m" + s + "\033[0m" }
func cyan(s string) string   { return "\033[36m" + s + "\033[0m" }
func white(s string) string  { return "\033[37m" + s + "\033[0m" }
func dim(s string) string    { return "\033[2m" + s + "\033[0m" }
