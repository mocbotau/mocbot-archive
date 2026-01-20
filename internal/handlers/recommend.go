package handlers

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mocbotau/mocbot-archive/internal/models"
	"github.com/mocbotau/mocbot-archive/internal/utils"
)

// GetRecommendedArtistsByUser retrieves recommended artists for a specific user.
func (h *Handler) GetRecommendedArtistsByUser(c *gin.Context) {
	userID := c.Param("userId")
	seedCount := c.Query("seedCount")

	seedCountInt, err := strconv.Atoi(seedCount)
	if err != nil || seedCountInt <= 0 {
		seedCountInt = utils.RecommendedTrackSeedCount
	}

	artists, err := h.db.GetRecentArtistsByUser(userID, seedCountInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve recommended artists"})
		return
	}

	c.JSON(200, gin.H{"recommended_artists": h.processRecommendedArtists(artists)})
}

// GetRecommendedArtistsByGuild retrieves recommended artists for a specific guild.
func (h *Handler) GetRecommendedArtistsByGuild(c *gin.Context) {
	guildID := c.Param("guildId")
	seedCount := c.Query("seedCount")

	seedCountInt, err := strconv.Atoi(seedCount)
	if err != nil || seedCountInt <= 0 {
		seedCountInt = utils.RecommendedTrackSeedCount
	}

	artists, err := h.db.GetRecentArtistsByGuild(guildID, seedCountInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve recommended artists"})
		return
	}

	c.JSON(200, gin.H{"recommended_artists": h.processRecommendedArtists(artists)})
}

// processRecommendedArtists processes a list of artists to calculate their weights and filter out bad artists.
func (h *Handler) processRecommendedArtists(artists []models.ArtistWithLatestPlay) []*models.RecommendedArtist {
	filteredArtists := h.filterBadArtists(artists)

	res := make([]*models.RecommendedArtist, 0, len(filteredArtists))

	for _, artist := range filteredArtists {
		res = append(res, &models.RecommendedArtist{
			Artist: artist.Artist,
			Weight: h.artistWeight(artist.LatestPlayTime, artist.Count, time.Now()),
		})
	}

	return res
}

// artistWeight calculates the weight of an artist based on recency and frequency of plays.
// Formula: weight = exp(-delta / tau) * log1p(count)
// where delta is the time since the latest play, tau is a time constant (6 hours), and count is the number of plays.
func (h *Handler) artistWeight(latest time.Time, count int, now time.Time) float64 {
	delta := now.Sub(latest).Seconds()
	recency := math.Exp(-delta / utils.Tau.Seconds())
	strength := math.Log1p(float64(count))

	w := recency * strength

	if w < 0.0001 {
		return 0.0001
	}

	return w
}

// filterBadArtists filters out artists that are considered "bad" based on a predefined list.
func (h *Handler) filterBadArtists(artists []models.ArtistWithLatestPlay) []*models.ArtistWithLatestPlay {
	var filtered []*models.ArtistWithLatestPlay

	for _, artist := range artists {
		if !slices.Contains(utils.BadArtists, strings.ToLower(artist.Artist)) {
			filtered = append(filtered, &artist)
		}
	}

	return filtered
}
