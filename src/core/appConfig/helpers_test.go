package appConfig

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestNormalizeConfigAllowsEmptyProxyAndDecodesBase64APIKey(t *testing.T) {
	rawAPIKey := "sk-test-key"
	encodedAPIKey := base64.StdEncoding.EncodeToString([]byte(rawAPIKey))

	config, err := NormalizeConfig(AssetxConfig{
		ProxyBaseURL: "",
		APIKey:       APIKeyBase64Prefix + encodedAPIKey,
	}, "test-config.json")
	if err != nil {
		t.Fatalf("NormalizeConfig returned error: %v", err)
	}

	if config.ProxyBaseURL != "" {
		t.Fatalf("Expected empty ProxyBaseURL, but got %q", config.ProxyBaseURL)
	}
	if config.APIKey != rawAPIKey {
		t.Fatalf("Expected decoded API key %q, but got %q", rawAPIKey, config.APIKey)
	}
}

func TestDecodeAPIKeyReturnsRawKeyWithoutPrefix(t *testing.T) {
	apiKey, err := DecodeAPIKey(" sk-raw-key ", "test-config.json")
	if err != nil {
		t.Fatalf("DecodeAPIKey returned error: %v", err)
	}

	if apiKey != "sk-raw-key" {
		t.Fatalf("Expected raw API key, but got %q", apiKey)
	}
}

func TestDecodeAPIKeyRejectsInvalidBase64(t *testing.T) {
	_, err := DecodeAPIKey(APIKeyBase64Prefix+"not base64", "test-config.json")
	if err == nil {
		t.Fatal("Expected invalid base64 error, but got nil")
	}

	if !strings.Contains(err.Error(), "invalid base64 api_key") {
		t.Fatalf("Expected invalid base64 error, but got %q", err.Error())
	}
}

func TestNormalizeConfigReturnsNilConfigOnError(t *testing.T) {
	config, err := NormalizeConfig(AssetxConfig{
		APIKey: APIKeyBase64Prefix + "not base64",
	}, "test-config.json")
	if err == nil {
		t.Fatal("Expected invalid base64 error, but got nil")
	}
	if config != nil {
		t.Fatalf("Expected nil config on error, but got %#v", config)
	}
}

func TestValidateConfigRejectsNilConfig(t *testing.T) {
	err := ValidateConfig(nil, "test-config.json")
	if err == nil {
		t.Fatal("Expected nil config error, but got nil")
	}

	if !strings.Contains(err.Error(), "config is nil") {
		t.Fatalf("Expected nil config error, but got %q", err.Error())
	}
}
