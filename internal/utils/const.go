package utils

import "time"

const (
	// DefaultLimit is the default maximum number of items to return in paginated requests.
	DefaultLimit = 50
	// MaxLimit is the maximum number of items that can be requested in paginated requests.
	MaxLimit = 500
	// RecommendedTracksCount is the number of recommended tracks to return.
	RecommendedTracksCount = 10
)

// DefaultSecretTTL is the default time-to-live for cached secrets.
const DefaultSecretTTL = 3 * time.Minute
