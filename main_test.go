package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binPath = "./testbin"

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Stderr.WriteString("failed to build test binary: " + string(out))
		os.Exit(1)
	}
	code := m.Run()
	os.Remove(binPath)
	os.Exit(code)
}

func run(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binPath, args...)
	var o, e strings.Builder
	cmd.Stdout = &o
	cmd.Stderr = &e
	exitCode = 0
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return o.String(), e.String(), exitCode
}

func TestCLI_NoArgs(t *testing.T) {
	_, stderr, code := run()
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("expected usage message in stderr, got:\n%s", stderr)
	}
}

func TestCLI_InvalidUnits(t *testing.T) {
	_, stderr, code := run("-units", "kelvin", "New York")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr, "invalid units") {
		t.Errorf("expected invalid units error, got:\n%s", stderr)
	}
}

func TestCLI_InvalidCompletion(t *testing.T) {
	_, stderr, code := run("-completion", "fish")
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr, "invalid completion shell") {
		t.Errorf("expected invalid shell error, got:\n%s", stderr)
	}
}

func TestCLI_CompletionBash(t *testing.T) {
	stdout, _, code := run("-completion", "bash")
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout, "complete -F _go_weather_demo") {
		t.Errorf("expected bash completion output, got:\n%s", stdout)
	}
}

func TestCLI_CompletionZsh(t *testing.T) {
	stdout, _, code := run("-completion", "zsh")
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout, "#compdef go_weather_demo") {
		t.Errorf("expected zsh completion output, got:\n%s", stdout)
	}
}
