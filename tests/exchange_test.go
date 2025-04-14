package tests

import (
	"testing"

	"github.com/astro/triarb/pkg/exchange"
)

func TestNewHitBTC(t *testing.T) {
	apiKey := "test-api-key"
	secretKey := "test-secret-key"

	hitbtc := exchange.NewHitBTC(apiKey, secretKey)

	if hitbtc.APIKey != apiKey {
		t.Errorf("Expected APIKey %s, got %s", apiKey, hitbtc.APIKey)
	}

	if hitbtc.SecretKey != secretKey {
		t.Errorf("Expected SecretKey %s, got %s", secretKey, hitbtc.SecretKey)
	}

	if hitbtc.Client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}
