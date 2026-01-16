package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// SecretManager handles secret file reading with TTL caching.
type SecretManager struct {
	cache *ttlcache.Cache[string, string]
}

// NewSecretManager creates a new SecretManager with the specified TTL.
// If ttl is 0, uses defaultSecretTTL (3 minutes).
func NewSecretManager(ttl time.Duration) *SecretManager {
	if ttl == 0 {
		ttl = DefaultSecretTTL
	}

	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](ttl),
	)
	go cache.Start()

	return &SecretManager{
		cache: cache,
	}
}

// GetSecretFromFile reads a secret from a specified file path, with a TTL cache.
func (sm *SecretManager) GetSecretFromFile(path string) (string, error) {
	item := sm.cache.Get(path)
	if item != nil && !item.IsExpired() {
		return item.Value(), nil
	}

	// #nosec G304 -- File path is controlled via environment variables
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read secret from file %s: %w", path, err)
	}

	secret := string(bytes)
	sm.cache.Set(path, secret, 0)

	return secret, nil
}

// Stop stops the secret manager's cache.
func (sm *SecretManager) Stop() {
	sm.cache.Stop()
}
