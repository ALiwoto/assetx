package appConfig

import (
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

func LoadConfig(configPath string) (AssetxConfig, error) {
	if strings.TrimSpace(configPath) == "" {
		defaultPath, err := DefaultConfigFilePath()
		if err != nil {
			return AssetxConfig{}, err
		}
		configPath = defaultPath
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return AssetxConfig{}, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	var config AssetxConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return AssetxConfig{}, fmt.Errorf("failed to parse config file %q as JSON: %w", configPath, err)
	}

	if err := ValidateConfig(config, configPath); err != nil {
		return AssetxConfig{}, err
	}

	return config, nil
}

func ValidateConfig(config AssetxConfig, configPath string) error {
	if strings.TrimSpace(config.ProxyBaseURL) == "" {
		return fmt.Errorf("config file %q is missing required field %q", configPath, "proxy_base_url")
	}
	if strings.TrimSpace(config.APIKey) == "" {
		return fmt.Errorf("config file %q is missing required field %q", configPath, "api_key")
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
