package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Well Known Environment Keys
const (
	EnvDbURL    = "B3SCALE_DB_URL"
	EnvLogLevel = "B3SCALE_LOG_LEVEL"
)

// Defaults
const (
	EnvDbURLDefault    = "postgres://postgres:postgres@localhost:5432/b3scale"
	EnvLogLevelDefault = "info"
)

// LoadEnv loads the environment from a file and
// updates the os.Environment
func LoadEnv(envfiles []string) {
	for _, filename := range envfiles {
		loadEnvFile(filename)
	}
}

// Internal load a single env file
func loadEnvFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		// We could not open the file - so we ignore this.
		return
	}
	defer f.Close()

	fmt.Println("using environment from:", filename)

	// Read lines and set env
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			continue // comment
		}
		tokens := strings.SplitN(l, "=", 2)
		if len(tokens) != 2 {
			continue // invalid
		}
		os.Setenv(
			strings.TrimSpace(tokens[0]),
			strings.TrimSpace(tokens[1]))
	}
}

// EnvOpt gets a configuration from the environment
// with a default fallback.
func EnvOpt(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
