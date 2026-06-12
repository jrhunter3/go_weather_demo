package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"go_weather_demo/output"
	"go_weather_demo/weather"
)

func main() {
	unitsFlag := flag.String("units", "imperial", "units: imperial or metric")
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "request timeout")
	jsonFlag := flag.Bool("json", false, "output as JSON")
	forecastFlag := flag.Int("forecast", 0, "show N-day forecast instead of current conditions")
	completionFlag := flag.String("completion", "", "generate shell completion script (bash or zsh)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <city> [city ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *completionFlag != "" {
		switch *completionFlag {
		case "bash":
			printBashCompletion()
		case "zsh":
			printZshCompletion()
		default:
			fmt.Fprintf(os.Stderr, "invalid completion shell %q: must be bash or zsh\n", *completionFlag)
			os.Exit(1)
		}
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var units output.Units
	switch *unitsFlag {
	case "imperial":
		units = output.Imperial
	case "metric":
		units = output.Metric
	default:
		fmt.Fprintf(os.Stderr, "invalid units %q: must be imperial or metric\n", *unitsFlag)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeoutFlag)
	defer cancel()

	start := time.Now()
	client := weather.NewClient()
	cities := flag.Args()

	if *forecastFlag > 0 {
		results := client.GetManyForecasts(ctx, cities, *forecastFlag)
		elapsed := time.Since(start)
		hasErr := false
		for _, r := range results {
			if r.Err != nil {
				hasErr = true
				break
			}
		}
		if *jsonFlag {
			output.PrintJSON(os.Stdout, results)
		} else {
			output.PrintForecast(os.Stdout, results, units, elapsed)
		}
		if hasErr {
			os.Exit(1)
		}
		return
	}

	results := client.GetMany(ctx, cities)
	elapsed := time.Since(start)

	if *jsonFlag {
		output.PrintJSON(os.Stdout, results)
	} else {
		output.Print(os.Stdout, results, units, elapsed)
	}

	for _, r := range results {
		if r.Err != nil {
			os.Exit(1)
			return
		}
	}
}

func printBashCompletion() {
	out := `_go_weather_demo() {
	local cur="${COMP_WORDS[COMP_CWORD]}"
	if [[ $cur == -* ]]; then
		COMPREPLY=($(compgen -W '-units -timeout -json -forecast -completion' -- "$cur"))
	fi
}
complete -F _go_weather_demo go_weather_demo
`
	os.Stdout.WriteString(out)
}

func printZshCompletion() {
	out := `#compdef go_weather_demo
_go_weather_demo() {
	_arguments \
		'--units[units: imperial or metric]' \
		'--timeout[request timeout]' \
		'--json[output as JSON]' \
		'--forecast[number of forecast days]' \
		'--completion[generate shell completion script]'
}
_go_weather_demo "$@"
`
	os.Stdout.WriteString(out)
}
