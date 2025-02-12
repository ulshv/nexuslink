package envutils

import (
	"fmt"
	"log"
	"os"

	"github.com/ulshv/nexuslink/internal/logger"
)

var envLogger = logger.NewSlogLogger("envutils")

func RequireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		envLogger.Error("environment variable is not set", "key", key)
		log.Fatal(fmt.Errorf("environment variable %s is not set", key))
	}
	return value
}

func OptionalEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		envLogger.Warn("environment variable is not set", "key", key)
		val = defaultValue
	}
	return val
}
