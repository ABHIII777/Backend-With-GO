package main

import (
	"os"
	"strings"
)

// loadDotEnv reads KEY=value lines from path (default ".env") and sets
// each key that is not already present in the environment. It is a
// no-op when the file is missing, so explicit environment variables
// (Docker, CI, fish/bash exports) always win.
//
// This keeps `go run ./cmd/server` working from any shell without
// requiring `source .env` (bash-only) or extra dependencies.
func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
