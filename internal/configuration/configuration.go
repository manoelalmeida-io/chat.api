package configuration

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func ResolveEnv(value string) string {
	start := "${"
	end := "}"

	if strings.HasPrefix(value, start) && strings.HasSuffix(value, end) {
		trimmed := strings.TrimPrefix(strings.TrimSuffix(value, end), start)
		parts := strings.SplitN(trimmed, ":", 2)
		key := parts[0]

		var defValue string

		if len(parts) > 1 {
			defValue = parts[1]
		}

		if envValue, exists := os.LookupEnv(key); exists {
			return envValue
		}

		return defValue
	}

	return value
}

func ResolveAllSettings(settings map[string]interface{}) map[string]interface{} {
	for key, value := range settings {
		switch v := value.(type) {
		case string:
			settings[key] = ResolveEnv(v)
		case map[string]interface{}:
			settings[key] = ResolveAllSettings(v)
		}
	}

	return settings
}

func ResolveAndUpdateAllSettings(config map[string]interface{}) {
	resolvedConfig := ResolveAllSettings(config)

	for key, value := range resolvedConfig {
		viper.Set(key, value)
	}
}
