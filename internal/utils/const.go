package utils

import "time"

const (
	// DefaultLimit is the default maximum number of items to return in paginated requests.
	DefaultLimit = 50
	// MaxLimit is the maximum number of items that can be requested in paginated requests.
	MaxLimit = 500
)

const (
	// MinTrackDurationMs is the minimum duration (in milliseconds) for a track to be considered valid.
	MinTrackDurationMs = 30 * 1000 // 30 seconds
	// MinTrackPercentagePlayed is the minimum percentage of a track that must be played for it to be considered valid.
	MinTrackPercentagePlayed = 0.15 // 15%
)

const (
	// RecommendedTrackSeedCount is the number of seed tracks to use for generating recommendations.
	RecommendedTrackSeedCount = 1000
	// Tau is the time constant used in the recommendation algorithm. Lower values favour more recent plays.
	Tau = 48 * time.Hour
)

// BadArtists is a set of artist names that are considered invalid or unhelpful for recommendations.
var BadArtists = []string{
	"unknown artist",
	"various artists",
	"various",
	"unknown",
	"n/a",
}

// DefaultSecretTTL is the default time-to-live for cached secrets.
const DefaultSecretTTL = 3 * time.Minute
