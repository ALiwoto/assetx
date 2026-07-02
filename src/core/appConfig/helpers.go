package appConfig

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func DefaultConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate the user home directory for default config path: %w", err)
	}

	return filepath.Join(homeDir, DefaultConfigRelativePath), nil
}

func LoadConfig(configPath string) (*AssetxConfig, error) {
	if strings.TrimSpace(configPath) == "" {
		defaultPath, err := DefaultConfigFilePath()
		if err != nil {
			return nil, err
		}
		configPath = defaultPath
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	var config = new(AssetxConfig)
	if err := json.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q as JSON: %w", configPath, err)
	}

	normalizedConfig, err := NormalizeConfig(config, configPath)
	if err != nil {
		return nil, err
	}

	if err := ValidateConfig(normalizedConfig, configPath); err != nil {
		return nil, err
	}

	return normalizedConfig, nil
}

func NormalizeConfig(config *AssetxConfig, configPath string) (*AssetxConfig, error) {
	config.ProxyBaseURL = strings.TrimSpace(config.ProxyBaseURL)

	decodedAPIKey, err := DecodeAPIKey(config.APIKey, configPath)
	if err != nil {
		return nil, err
	}
	config.APIKey = decodedAPIKey

	return config, nil
}

func DecodeAPIKey(apiKey string, configPath string) (string, error) {
	trimmedAPIKey := strings.TrimSpace(apiKey)
	if !strings.HasPrefix(trimmedAPIKey, APIKeyBase64Prefix) {
		return trimmedAPIKey, nil
	}

	encodedAPIKey := strings.TrimSpace(strings.TrimPrefix(trimmedAPIKey, APIKeyBase64Prefix))
	decodedAPIKeyBytes, err := base64.StdEncoding.DecodeString(encodedAPIKey)
	if err != nil {
		return "", fmt.Errorf("config file %q has invalid base64 api_key after %q prefix: %w", configPath, APIKeyBase64Prefix, err)
	}

	decodedAPIKey := strings.TrimSpace(string(decodedAPIKeyBytes))
	if decodedAPIKey == "" {
		return "", fmt.Errorf("config file %q has empty api_key after %q base64 decoding", configPath, APIKeyBase64Prefix)
	}

	return decodedAPIKey, nil
}

func ValidateConfig(config *AssetxConfig, configPath string) error {
	if config == nil {
		return fmt.Errorf("config file %q could not be validated because config is nil", configPath)
	}

	if strings.TrimSpace(config.APIKey) == "" {
		return fmt.Errorf("config file %q is missing required field %q", configPath, "api_key")
	}

	if strings.TrimSpace(config.ProxyBaseURL) == "" {
		return nil
	}

	parsedURL, err := url.ParseRequestURI(config.ProxyBaseURL)
	if err != nil {
		return fmt.Errorf("config file %q has invalid proxy_base_url %q: %w", configPath, config.ProxyBaseURL, err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("config file %q has invalid proxy_base_url %q: expected http or https URL", configPath, config.ProxyBaseURL)
	}

	return nil
}
